// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

string constant TAG = "PriceFeed";

contract PriceFeed is Initializable, ContextUpgradeable {
    using SafeMathUpgradeable for uint256;

    address public roleManager;
    mapping(bytes32 => uint256) public prices;
    mapping(bytes32 => uint256) public proposedPrices;
    mapping(bytes32 => bool) public proposedStatus;
    uint256 public DEVIATION_THRESHOLD;
    bytes32 public constant PROPOSE_PRICE = keccak256("PROPOSE_PRICE");
    bytes32 public constant VALIDATE_PRICE = keccak256("VALIDATE_PRICE");
    mapping(bytes32 => uint256) public proposedCount;

    event SetPrice(bytes32 symbol, uint256 price, uint256 blockTimestamp);
    event ProposePrice(
        bytes32[] symbol,
        uint256[] price,
        uint256 blockTimestamp
    );
    event DeviationExceeded(
        bytes32 symbol,
        uint256 validatedPrice,
        uint256 proposedPrice,
        uint256 blockTimestamp,
        uint256 iterationCount
    );
    event IterationExceeded(bytes32 symbol);

    error ArrayLengthMismatched(string contractTAG);
    error PriceNotSet(string contractTAG);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initilize the contract.
     * @param roleManagerAddress - Address of the roleManager contract.
     **/
    function initialize(address roleManagerAddress) public initializer {
        roleManager = roleManagerAddress;
    }

    /**
     * @notice Configures the deviation threshold between proposed and validated prices
     * This function can only be accessed by the account with role SET_ORACLE_METHOD.
     * @param _threshold - threshold deviation
     **/
    function setDeviationThreshold(uint256 _threshold) external {
        DEVIATION_THRESHOLD = _threshold;
    }

    /**
     * @notice Proposes the prices of assets
     * This function can only be accessed by the account with role PROPOSE_PRICE.
     * @param _symbols - list of pool ids of the assets
     * @param _prices - prices of the asset
     **/
    function proposePrices(
        bytes32[] memory _symbols,
        uint256[] memory _prices
    ) external {
        if (_symbols.length != _prices.length)
            revert ArrayLengthMismatched(TAG);
        for (uint256 i; i < _symbols.length; i++) {
            proposedPrices[_symbols[i]] = _prices[i];
            proposedStatus[_symbols[i]] = true;
        }
        emit ProposePrice(_symbols, _prices, block.timestamp);
    }

    /**
     * @notice Validates the prices of assets
     * This function can only be accessed by the account with role VALIDATE_PRICE.
     * @param _symbols - list of pool ids of the assets
     * @param _prices - prices of the asset
     **/
    function validatePrices(
        bytes32[] memory _symbols,
        uint256[] memory _prices
    ) external {
        uint256 minPrice;
        uint256 maxPrice;
        uint256 deviation;
        if (_symbols.length != _prices.length)
            revert ArrayLengthMismatched(TAG);
        for (uint256 i; i < _symbols.length; i++) {
            if (_prices[i] > proposedPrices[_symbols[i]]) {
                minPrice = proposedPrices[_symbols[i]];
                maxPrice = _prices[i];
            } else {
                minPrice = _prices[i];
                maxPrice = proposedPrices[_symbols[i]];
            }
            deviation = ((maxPrice.sub(minPrice)).mul(1e18));
            deviation = deviation.div(proposedPrices[_symbols[i]]);
            if (deviation > DEVIATION_THRESHOLD) {
                proposedCount[_symbols[i]]++;
                emit DeviationExceeded(
                    _symbols[i],
                    _prices[i],
                    proposedPrices[_symbols[i]],
                    block.timestamp,
                    proposedCount[_symbols[i]]
                );
                proposedStatus[_symbols[i]] = false;
                continue;
            }
            prices[_symbols[i]] = _prices[i];
            proposedStatus[_symbols[i]] = false;
            proposedCount[_symbols[i]] = 0;
            emit SetPrice(_symbols[i], _prices[i], block.timestamp);
        }
    }

    /**
     * @notice Returns price of assets
     * @param _symbols - list of the asset's pool id
     **/
    function getPrices(
        bytes32[] memory _symbols
    ) external view returns (uint256[] memory results) {
        results = new uint256[](_symbols.length);
        for (uint256 i; i < _symbols.length; i++) {
            uint256 price = prices[_symbols[i]];
            if (price <= 0) revert PriceNotSet(TAG);
            results[i] = price;
        }
    }

    /**
     * @notice Returns price of asset
     * @param _symbol - asset's pool id
     **/
    function getPrice(bytes32 _symbol) external view returns (uint256 price) {
        price = prices[_symbol];
        if (price <= 0) revert PriceNotSet(TAG);
    }
}

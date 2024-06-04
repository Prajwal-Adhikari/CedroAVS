// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "./ERC20Mock.sol";
import "../../interfaces/IOracle.sol";

contract AaveMock {
    using SafeMath for uint256;
    mapping(address => address) public aToken;
    mapping(address => address) public dToken;
    mapping(address => address[]) public userCollateralList;
    mapping(address => address[]) public userDebtList;
    mapping(address => mapping(address => bool)) public hasToken;
    mapping(address => mapping(address => bool)) public hasDebt;
    mapping(address => bytes32) public assetKey;
    address[] public assets;
    mapping(address => bool) public isIncluded;
    address public oracle;

    constructor(address _oracle) {
        oracle = _oracle;
    }

    function setAssetKey(address _asset, bytes32 _key) public {
        assetKey[_asset] = _key;
    }

    function setAToken(address _asset, address _aToken) public {
        aToken[_asset] = _aToken;
    }

    function setDToken(address _asset, address _dToken) public {
        dToken[_asset] = _dToken;
    }

    function setAssetsList(address[] memory _assets) public {
        uint256 len = _assets.length;
        for (uint256 i; i < len; ) {
            if (!isIncluded[_assets[i]]) {
                assets.push(_assets[i]);
            }
            unchecked {
                i++;
            }
        }
    }

    function supply(
        address asset,
        uint256 amount,
        address onBehalfOf,
        uint16 referralCode
    ) public virtual {
        ERC20Mock(aToken[asset]).mint(onBehalfOf, amount);
        _pay(asset, amount);
    }

    function withdraw(
        address asset,
        uint256 amount,
        address to
    ) public virtual returns (uint256) {
        ERC20Mock(aToken[asset]).burn(msg.sender, amount);
        IERC20(asset).transfer(to, amount);
        checkHealthFactor(msg.sender);

        return amount;
    }

    function borrow(
        address asset,
        uint256 amount,
        uint256 interestRateMode,
        uint16 referralCode,
        address onBehalfOf
    ) public {
        ERC20Mock(dToken[asset]).mint(msg.sender, amount);
        IERC20(asset).transfer(onBehalfOf, amount);
        checkHealthFactor(msg.sender);
    }

    function repay(
        address asset,
        uint256 amount,
        uint256 interestRateMode,
        address onBehalfOf
    ) public returns (uint256) {
        ERC20Mock(dToken[asset]).burn(onBehalfOf, amount);
        _pay(asset, amount);
        return amount;
    }

    function repayWithATokens(
        address asset,
        uint256 amount,
        uint256 interestRateMode
    ) public returns (uint256) {
        ERC20Mock(dToken[asset]).burn(msg.sender, amount);
        ERC20Mock(aToken[asset]).burn(msg.sender, amount);
    }

    function _pay(address _token, uint256 _amount) internal {
        if (_token == address(0)) {
            if (msg.value < _amount) revert("amount not sufficient");
        } else {
            IERC20(_token).transferFrom(msg.sender, address(this), _amount);
        }
    }

    function testPayload(
        uint8 functionType,
        string memory user,
        address id,
        uint256 remoteAmount
    ) external pure returns (bytes memory) {
        bytes memory payload = abi.encode(
            functionType,
            stringToBytes32(user),
            addressToBytes32(id),
            remoteAmount
        );

        return payload;
    }

    function stringToBytes32(string memory source)
        public
        pure
        returns (bytes32 result)
    {
        assembly {
            result := mload(add(source, 32))
        }
    }

    function checkHealthFactor(address user) internal {
        uint256 collateral;
        uint256 debt;
        uint256 len = assets.length;
        for (uint256 i; i < len; ) {
            collateral += ERC20Mock(aToken[assets[i]])
                .balanceOf(user)
                .mul(IOracle(oracle).getSinglePrice(assetKey[assets[i]]))
                .div(1e18);
            unchecked {
                i++;
            }
        }

        for (uint256 i; i < len; ) {
            debt += ERC20Mock(dToken[assets[i]])
                .balanceOf(user)
                .mul(IOracle(oracle).getSinglePrice(assetKey[assets[i]]))
                .div(1e18);

            unchecked {
                i++;
            }
        }
        require(
            debt == 0 || debt.mul(1e18).div(collateral) <= 5e17,
            "Health Factor exceeded"
        );
    }

    function addressToBytes32(address addr) public pure returns (bytes32) {
        bytes32 result;
        assembly {
            mstore(0, addr)
            result := mload(0)
        }
        return result;
    }
}

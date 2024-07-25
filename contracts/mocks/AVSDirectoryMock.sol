// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IAVSDirectory, ISignatureUtils} from "../../contracts/core/eigenContracts/interfaces/IAVSDirectory.sol";
import {IDelegationManager} from "../../contracts/core/eigenContracts/interfaces/IDelegationManager.sol";
import {EIP1271SignatureUtils} from "../../contracts/core/eigenContracts/libraries/EIP1271SignatureUtils.sol";


contract AVSDirectoryMock is IAVSDirectory {
     IDelegationManager public immutable delegation;
      // @dev Chain ID at the time of contract deployment
    uint256 internal immutable ORIGINAL_CHAIN_ID;
    /// @notice Mapping: AVS => operator => enum of operator status to the AVS
    mapping(address => mapping(address => OperatorAVSRegistrationStatus)) public avsOperatorStatus;

    /// @notice Mapping: operator => 32-byte salt => whether or not the salt has already been used by the operator.
    /// @dev Salt is used in the `registerOperatorToAVS` function.
    mapping(address => mapping(bytes32 => bool)) public operatorSaltIsSpent;

     constructor( IDelegationManager _delegation) {
        delegation = _delegation;
        ORIGINAL_CHAIN_ID = block.chainid;
    }
    /**
     * @notice Called by an avs to register an operator with the avs.
     * @param operator The address of the operator to register.
     * @param operatorSignature The signature, salt, and expiry of the operator's signature.
     */
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external  {

        require(
            operatorSignature.expiry >= block.timestamp,
            "AVSDirectory.registerOperatorToAVS: operator signature expired"
        );
        require(
            avsOperatorStatus[msg.sender][operator] != OperatorAVSRegistrationStatus.REGISTERED,
            "AVSDirectory.registerOperatorToAVS: operator already registered"
        );
        require(
            !operatorSaltIsSpent[operator][operatorSignature.salt],
            "AVSDirectory.registerOperatorToAVS: salt already spent"
        );
        require(
            delegation.isOperator(operator),
            "AVSDirectory.registerOperatorToAVS: operator not registered to EigenLayer yet");

        // Calculate the digest hash
        bytes32 operatorRegistrationDigestHash =  calculateOperatorAVSRegistrationDigestHash({
            operator: operator,
            avs: msg.sender,
            salt: operatorSignature.salt,
            expiry: operatorSignature.expiry
        });
     
     

        // // Check that the signature is valid
        EIP1271SignatureUtils.checkSignature_EIP1271(
            operator,
            operatorRegistrationDigestHash,
            operatorSignature.signature
        );
        
        // Set the operator as registered
        avsOperatorStatus[msg.sender][operator] = OperatorAVSRegistrationStatus.REGISTERED;

        // Mark the salt as spent
        operatorSaltIsSpent[operator][operatorSignature.salt] = true;

        emit OperatorAVSRegistrationStatusUpdated(operator, msg.sender, OperatorAVSRegistrationStatus.REGISTERED);
    }

    /**
     * @notice Called by an avs to deregister an operator with the avs.
     * @param operator The address of the operator to deregister.
     */
    function deregisterOperatorFromAVS(address operator) external {}

    /**
     * @notice Called by an AVS to emit an `AVSMetadataURIUpdated` event indicating the information has updated.
     * @param metadataURI The URI for metadata associated with an AVS
     * @dev Note that the `metadataURI` is *never stored * and is only emitted in the `AVSMetadataURIUpdated` event
     */
    function updateAVSMetadataURI(string calldata metadataURI) external {}

    /**
     * @notice Returns whether or not the salt has already been used by the operator.
     * @dev Salts is used in the `registerOperatorToAVS` function.
     */
    // function operatorSaltIsSpent(address operator, bytes32 salt) external view returns (bool) {}

    /**
     * @notice Calculates the digest hash to be signed by an operator to register with an AVS
     * @param operator The account registering as an operator
     * @param avs The AVS the operator is registering to
     * @param salt A unique and single use value associated with the approver signature.
     * @param expiry Time after which the approver's signature becomes invalid
     */
    function calculateOperatorAVSRegistrationDigestHash(
        address operator,
        address avs,
        bytes32 salt,
        uint256 expiry
    ) public view returns (bytes32) {
         bytes32 structHash = keccak256(
         abi.encode(OPERATOR_AVS_REGISTRATION_TYPEHASH(), operator, avs, salt, expiry)
        );
        // calculate the digest hash
        bytes32 digestHash = keccak256(
            abi.encodePacked("\x19\x01", delegation.domainSeparator(), structHash)
        );
        return digestHash;
    }

    /// @notice The EIP-712 typehash for the Registration struct used by the contract
    function OPERATOR_AVS_REGISTRATION_TYPEHASH() public pure returns (bytes32) {
        return keccak256("OperatorAVSRegistration(address operator,address avs,bytes32 salt,uint256 expiry)");
    }

} 
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mevboostrelaylist

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	Uri         string
	Operator    string
	IsMandatory bool
	Description string
}

// ApiMetaData contains all meta data concerning the Api contract.
var ApiMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"RelayAdded\",\"inputs\":[{\"name\":\"uri_hash\",\"type\":\"string\",\"indexed\":true},{\"name\":\"relay\",\"type\":\"tuple\",\"components\":[{\"name\":\"uri\",\"type\":\"string\"},{\"name\":\"operator\",\"type\":\"string\"},{\"name\":\"is_mandatory\",\"type\":\"bool\"},{\"name\":\"description\",\"type\":\"string\"}],\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RelayRemoved\",\"inputs\":[{\"name\":\"uri_hash\",\"type\":\"string\",\"indexed\":true},{\"name\":\"uri\",\"type\":\"string\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"AllowedListUpdated\",\"inputs\":[{\"name\":\"allowed_list_version\",\"type\":\"uint256\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"OwnerChanged\",\"inputs\":[{\"name\":\"new_owner\",\"type\":\"address\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"ManagerChanged\",\"inputs\":[{\"name\":\"new_manager\",\"type\":\"address\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"ERC20Recovered\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_relays_amount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_manager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_relays\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"components\":[{\"name\":\"uri\",\"type\":\"string\"},{\"name\":\"operator\",\"type\":\"string\"},{\"name\":\"is_mandatory\",\"type\":\"bool\"},{\"name\":\"description\",\"type\":\"string\"}]}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_relay_by_uri\",\"inputs\":[{\"name\":\"relay_uri\",\"type\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"components\":[{\"name\":\"uri\",\"type\":\"string\"},{\"name\":\"operator\",\"type\":\"string\"},{\"name\":\"is_mandatory\",\"type\":\"bool\"},{\"name\":\"description\",\"type\":\"string\"}]}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_allowed_list_version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"add_relay\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\"},{\"name\":\"operator\",\"type\":\"string\"},{\"name\":\"is_mandatory\",\"type\":\"bool\"},{\"name\":\"description\",\"type\":\"string\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"remove_relay\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"change_owner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_manager\",\"inputs\":[{\"name\":\"manager\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"dismiss_manager\",\"inputs\":[],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"recover_erc20\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"}]",
	Bin: "0x60206115ff6000396000518060a01c6115fa57604052346115fa5760405161007e5760126060527f7a65726f206f776e65722061646472657373000000000000000000000000000060805260605060605180608001601f826000031636823750506308c379a06020526020604052601f19601f6060510116604401603cfd5b60405160005561156361009661000039611563610000f36003361161000c57611181565b60003560e01c346115515763312c3165811861003657600436186115515760025460405260206040f35b630ac298dc811861005557600436186115515760005460405260206040f35b639e4a0fc4811861007457600436186115515760015460405260206040f35b6304e469ea81186102545760043618611551576020806040528060400160006002548083528060051b6000826028811161155157801561024057905b828160051b60208801015260648102600301836020880101608080825280820183548082526001850160208301600083601f0160051c6020811161155157801561010c57905b808401548160051b8401526001018181186100f6575b50505050508051806020830101601f82600003163682375050601f19601f825160200101169050810190508060208301526021830181830181548082526001830160208301600083601f0160051c6020811161155157801561018057905b808401548160051b84015260010181811861016a575b50505050508051806020830101601f82600003163682375050601f19601f825160200101169050905081019050604283015460408301528060608301526043830181830181548082526001830160208301600083601f0160051c6020811161155157801561020057905b808401548160051b8401526001018181186101ea575b50505050508051806020830101601f82600003163682375050601f19601f82516020010116905090508101905090509050830192506001018181186100b0575b505082016020019150509050810190506040f35b63f5f33c7b81186104b457604436106115515760043560040161040081351161155157803580611120526020820181816111403750505061112051806040528060608261114060045afa50506102ab611560611187565b61156051611540526115405119610322576015611560527f6e6f2072656c61792077697468207468652055524900000000000000000000006115805261156050611560518061158001601f826000031636823750506308c379a061152052602061154052601f19601f61156051011660440161153cfd5b6020806115605260646115405160025481101561155157026003018161156001608080825280820183548082526001850160208301600083601f0160051c6020811161155157801561038657905b808401548160051b840152600101818118610370575b50505050508051806020830101601f82600003163682375050601f19601f825160200101169050810190508060208301526021830181830181548082526001830160208301600083601f0160051c602081116115515780156103fa57905b808401548160051b8401526001018181186103e4575b50505050508051806020830101601f82600003163682375050601f19601f825160200101169050905081019050604283015460408301528060608301526043830181830181548082526001830160208301600083601f0160051c6020811161155157801561047a57905b808401548160051b840152600101818118610464575b50505050508051806020830101601f82600003163682375050601f19601f8251602001011690509050810190509050905081019050611560f35b6376650ad381186104d4576004361861155157610fa35460405260206040f35b632e21ecef81186109805760e43610611551576004356004016104008135116115515780358061112052602082018181611140375050506024356004016104008135116115515780358061154052602082018181611560375050506044358060011c611551576119605260643560040161040081351161155157803580611980526020820181816119a03750505061056a6112ed565b6000611da052611da08051602082012090506111205161114020186105ef57601b6121c0527f72656c617920555249206d757374206e6f7420626520656d70747900000000006121e0526121c0506121c051806121e001601f826000031636823750506308c379a06121805260206121a052601f19601f6121c051011660440161219cfd5b6027600254111561066057601c611da0527f616c7265616479206d6178206e756d626572206f662072656c61797300000000611dc052611da050611da05180611dc001601f826000031636823750506308c379a0611d60526020611d8052601f19601f611da0510116604401611d7cfd5b61112051806040528060608261114060045afa5050610680611dc0611187565b611dc051611da052611da051191561071d576021611dc0527f72656c61792077697468207468652055524920616c7265616479206578697374611de0527f7300000000000000000000000000000000000000000000000000000000000000611e0052611dc050611dc05180611de001601f826000031636823750506308c379a0611d80526020611da052601f19601f611dc0510116604401611d9cfd5b6111205180611dc05280611de08261114060045afa505061154051806121e052806122008261156060045afa5050611960516126005261198051806126205280612640826119a060045afa505060025460278111611551576001810160025560648102600301611dc05180825560018201600082601f0160051c602081116115515780156107bf57905b8060051b611de00151818401556001018181186107a7575b505050506121e05180602183015560016021830101600082601f0160051c6020811161155157801561080557905b8060051b6122000151818401556001018181186107ed575b505050506126005160428201556126205180604383015560016043830101600082601f0160051c6020811161155157801561085457905b8060051b61264001518184015560010181811861083c575b505050505050610862611374565b61112051611140207feee5faa84d45af657ab405cdbf2c6a8a6d466e83fa694a358fee5ff84431d0bf602080612a405280612a40016080808252808201611dc05180825260208201818183611de060045afa5050508051806020830101601f82600003163682375050601f19601f825160200101169050810190508060208301528082016121e0518082526020820181818361220060045afa5050508051806020830101601f82600003163682375050601f19601f82516020010116905081019050612600516040830152806060830152808201612620518082526020820181818361264060045afa5050508051806020830101601f82600003163682375050601f19601f82516020010116905081019050905081019050612a40a2005b63f5a70a808118610ca45760443610611551576004356004016104008135116115515780358061112052602082018181611140375050506109bf6112ed565b600061154052611540805160208201209050611120516111402018610a4457601b611960527f72656c617920555249206d757374206e6f7420626520656d70747900000000006119805261196050611960518061198001601f826000031636823750506308c379a061192052602061194052601f19601f61196051011660440161193cfd5b6002546115405261112051806040528060608261114060045afa5050610a6b611580611187565b6115805161156052611540516115605110610ae6576015611580527f6e6f2072656c61792077697468207468652055524900000000000000000000006115a0526115805061158051806115a001601f826000031636823750506308c379a061154052602061156052601f19601f61158051011660440161155cfd5b61154051600181038181116115515790506115605114610c1257606461156051600254811015611551570260030160646115405160018103818111611551579050600254811015611551570260030180548083556001820160018401600083601f0160051c60208111611551578015610b6e57905b8084015481840155600101818118610b5b575b50505050506021810180548060218501556001820160016021860101600083601f0160051c60208111611551578015610bb657905b8084015481840155600101818118610ba3575b505050505050604281015460428301556043810180548060438501556001820160016043860101600083601f0160051c60208111611551578015610c0957905b8084015481840155600101818118610bf6575b50505050505050505b6001600254801561155157038060025550610c2b611374565b61112051611140207fef756634af7ee7210f786ec0f91930afa63fda84d9ff6493ae681c332055dadb602080611580528061158001611120518082526020820181818361114060045afa5050508051806020830101601f82600003163682375050601f19601f82516020010116905081019050611580a2005b63253c8bd48118610dca5760243618611551576004358060a01c61155157608052610ccd6113ba565b608051610d3157601260a0527f7a65726f206f776e65722061646472657373000000000000000000000000000060c05260a05060a0518060c001601f826000031636823750506308c379a06060526020608052601f19601f60a0510116604401607cfd5b60005460805118610d9957600a60a0527f73616d65206f776e65720000000000000000000000000000000000000000000060c05260a05060a0518060c001601f826000031636823750506308c379a06060526020608052601f19601f60a0510116604401607cfd5b6080516000556080517fa2ea9883a321a3e97b8266c2b078bfeec6d50c711ed71f874a90d500ae2eaf36600060a0a2005b639aece83e8118610ef05760243618611551576004358060a01c61155157608052610df36113ba565b608051610e5757601460a0527f7a65726f206d616e61676572206164647265737300000000000000000000000060c05260a05060a0518060c001601f826000031636823750506308c379a06060526020608052601f19601f60a0510116604401607cfd5b60015460805118610ebf57600c60a0527f73616d65206d616e61676572000000000000000000000000000000000000000060c05260a05060a0518060c001601f826000031636823750506308c379a06060526020608052601f19601f60a0510116604401607cfd5b6080516001556080517f198db6e425fb8aafd1823c6ca50be2d51e5764571a5ae0f0f21c6812e45def0b600060a0a2005b63417a02b48118610f9e576004361861155157610f0b6113ba565b600154610f6f57600e6080527f6e6f206d616e616765722073657400000000000000000000000000000000000060a0526080506080518060a001601f826000031636823750506308c379a06040526020606052601f19601f6080510116604401605cfd5b600060015560007f198db6e425fb8aafd1823c6ca50be2d51e5764571a5ae0f0f21c6812e45def0b60006080a2005b63edd885b4811861117f5760643618611551576004358060a01c611551576101e0526044358060a01c6115515761020052610fd76113ba565b6101e051611045576012610220527f7a65726f20746f6b656e206164647265737300000000000000000000000000006102405261022050610220518061024001601f826000031636823750506308c379a06101e052602061020052601f19601f6102205101166044016101fcfd5b610200516110b3576016610220527f7a65726f20726563697069656e742061646472657373000000000000000000006102405261022050610220518061024001601f826000031636823750506308c379a06101e052602061020052601f19601f6102205101166044016101fcfd5b6101e0513b611122576011610220527f656f6120746f6b656e20616464726573730000000000000000000000000000006102405261022050610220518061024001601f826000031636823750506308c379a06101e052602061020052601f19601f6102205101166044016101fcfd5b6024351561117d576101e05160405261020051606052602435608052611146611423565b610200516101e0517f8619312ed4eff1cf9f0116e6db2f49d9570a86f0350d1c5ad1bd0f7b0cf9e132602435610220526020610220a35b005b505b60006000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff610460526000610480526000600254602881116115515780156112e257905b606481026003018054806104a05260018201600082601f0160051c6020811161155157801561120a57905b808301548160051b6104c001526001018181186111f2575b50505050602181018054806108c05260018201600082601f0160051c6020811161155157801561124e57905b808301548160051b6108e00152600101818118611236575b50505050506042810154610ce05260438101805480610d005260018201600082601f0160051c6020811161155157801561129c57905b808301548160051b610d200152600101818118611284575b5050505050506040516060206104a0516104c020186112c25761048051610460526112e2565b6104805160018101818110611551579050610480526001018181186111c7575b505061046051815250565b60005433186112fd576001611311565b600154331861130e57331515611311565b60005b61137257601f6040527f6d73672e73656e646572206e6f74206f776e6572206f72206d616e616765720060605260405060405180606001601f826000031636823750506308c379a06000526020602052601f19601f6040510116604401601cfd5b565b610fa35460018101818110611551579050604052604051610fa3556040517f49f5627aa055ec3fcd474f99c8b7799b798c04af7b9f215305512c867e5a183960006060a2565b6000543318156114215760146040527f6d73672e73656e646572206e6f74206f776e657200000000000000000000000060605260405060405180606001601f826000031636823750506308c379a06000526020602052601f19601f6040510116604401601cfd5b565b6000600460e0527fa9059cbb000000000000000000000000000000000000000000000000000000006101005260e08051602082018361014001815181525050808301925050506060518161014001526020810190506080518161014001526020810190508061012052610120505060206101c06101205161014060006040515af16114b3573d600060003e3d6000fd5b3d602081183d60201002186101a0526101a080518060a05260208201805160c05250505060a0511561154f5760c05160a05160200360031b1c61154f57601560e0527f6572633230207472616e73666572206661696c656400000000000000000000006101005260e05060e0518061010001601f826000031636823750506308c379a060a052602060c052601f19601f60e051011660440160bcfd5b565b600080fda165767970657283000306000b005b600080fd",
}

// ApiABI is the input ABI used to generate the binding from.
// Deprecated: Use ApiMetaData.ABI instead.
var ApiABI = ApiMetaData.ABI

// ApiBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ApiMetaData.Bin instead.
var ApiBin = ApiMetaData.Bin

// DeployApi deploys a new Ethereum contract, binding an instance of Api to it.
func DeployApi(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address) (common.Address, *types.Transaction, *Api, error) {
	parsed, err := ApiMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ApiBin), backend, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Api{ApiCaller: ApiCaller{contract: contract}, ApiTransactor: ApiTransactor{contract: contract}, ApiFilterer: ApiFilterer{contract: contract}}, nil
}

// Api is an auto generated Go binding around an Ethereum contract.
type Api struct {
	ApiCaller     // Read-only binding to the contract
	ApiTransactor // Write-only binding to the contract
	ApiFilterer   // Log filterer for contract events
}

// ApiCaller is an auto generated read-only Go binding around an Ethereum contract.
type ApiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ApiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ApiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ApiSession struct {
	Contract     *Api              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ApiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ApiCallerSession struct {
	Contract *ApiCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ApiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ApiTransactorSession struct {
	Contract     *ApiTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ApiRaw is an auto generated low-level Go binding around an Ethereum contract.
type ApiRaw struct {
	Contract *Api // Generic contract binding to access the raw methods on
}

// ApiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ApiCallerRaw struct {
	Contract *ApiCaller // Generic read-only contract binding to access the raw methods on
}

// ApiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ApiTransactorRaw struct {
	Contract *ApiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewApi creates a new instance of Api, bound to a specific deployed contract.
func NewApi(address common.Address, backend bind.ContractBackend) (*Api, error) {
	contract, err := bindApi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Api{ApiCaller: ApiCaller{contract: contract}, ApiTransactor: ApiTransactor{contract: contract}, ApiFilterer: ApiFilterer{contract: contract}}, nil
}

// NewApiCaller creates a new read-only instance of Api, bound to a specific deployed contract.
func NewApiCaller(address common.Address, caller bind.ContractCaller) (*ApiCaller, error) {
	contract, err := bindApi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ApiCaller{contract: contract}, nil
}

// NewApiTransactor creates a new write-only instance of Api, bound to a specific deployed contract.
func NewApiTransactor(address common.Address, transactor bind.ContractTransactor) (*ApiTransactor, error) {
	contract, err := bindApi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ApiTransactor{contract: contract}, nil
}

// NewApiFilterer creates a new log filterer instance of Api, bound to a specific deployed contract.
func NewApiFilterer(address common.Address, filterer bind.ContractFilterer) (*ApiFilterer, error) {
	contract, err := bindApi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ApiFilterer{contract: contract}, nil
}

// bindApi binds a generic wrapper to an already deployed contract.
func bindApi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ApiMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Api *ApiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Api.Contract.ApiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Api *ApiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Api.Contract.ApiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Api *ApiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Api.Contract.ApiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Api *ApiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Api.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Api *ApiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Api.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Api *ApiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Api.Contract.contract.Transact(opts, method, params...)
}

// GetAllowedListVersion is a free data retrieval call binding the contract method 0x76650ad3.
//
// Solidity: function get_allowed_list_version() view returns(uint256)
func (_Api *ApiCaller) GetAllowedListVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_allowed_list_version")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetAllowedListVersion is a free data retrieval call binding the contract method 0x76650ad3.
//
// Solidity: function get_allowed_list_version() view returns(uint256)
func (_Api *ApiSession) GetAllowedListVersion() (*big.Int, error) {
	return _Api.Contract.GetAllowedListVersion(&_Api.CallOpts)
}

// GetAllowedListVersion is a free data retrieval call binding the contract method 0x76650ad3.
//
// Solidity: function get_allowed_list_version() view returns(uint256)
func (_Api *ApiCallerSession) GetAllowedListVersion() (*big.Int, error) {
	return _Api.Contract.GetAllowedListVersion(&_Api.CallOpts)
}

// GetManager is a free data retrieval call binding the contract method 0x9e4a0fc4.
//
// Solidity: function get_manager() view returns(address)
func (_Api *ApiCaller) GetManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_manager")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetManager is a free data retrieval call binding the contract method 0x9e4a0fc4.
//
// Solidity: function get_manager() view returns(address)
func (_Api *ApiSession) GetManager() (common.Address, error) {
	return _Api.Contract.GetManager(&_Api.CallOpts)
}

// GetManager is a free data retrieval call binding the contract method 0x9e4a0fc4.
//
// Solidity: function get_manager() view returns(address)
func (_Api *ApiCallerSession) GetManager() (common.Address, error) {
	return _Api.Contract.GetManager(&_Api.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x0ac298dc.
//
// Solidity: function get_owner() view returns(address)
func (_Api *ApiCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetOwner is a free data retrieval call binding the contract method 0x0ac298dc.
//
// Solidity: function get_owner() view returns(address)
func (_Api *ApiSession) GetOwner() (common.Address, error) {
	return _Api.Contract.GetOwner(&_Api.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x0ac298dc.
//
// Solidity: function get_owner() view returns(address)
func (_Api *ApiCallerSession) GetOwner() (common.Address, error) {
	return _Api.Contract.GetOwner(&_Api.CallOpts)
}

// GetRelayByUri is a free data retrieval call binding the contract method 0xf5f33c7b.
//
// Solidity: function get_relay_by_uri(string relay_uri) view returns((string,string,bool,string))
func (_Api *ApiCaller) GetRelayByUri(opts *bind.CallOpts, relay_uri string) (Struct0, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_relay_by_uri", relay_uri)
	if err != nil {
		return *new(Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct0)).(*Struct0)

	return out0, err
}

// GetRelayByUri is a free data retrieval call binding the contract method 0xf5f33c7b.
//
// Solidity: function get_relay_by_uri(string relay_uri) view returns((string,string,bool,string))
func (_Api *ApiSession) GetRelayByUri(relay_uri string) (Struct0, error) {
	return _Api.Contract.GetRelayByUri(&_Api.CallOpts, relay_uri)
}

// GetRelayByUri is a free data retrieval call binding the contract method 0xf5f33c7b.
//
// Solidity: function get_relay_by_uri(string relay_uri) view returns((string,string,bool,string))
func (_Api *ApiCallerSession) GetRelayByUri(relay_uri string) (Struct0, error) {
	return _Api.Contract.GetRelayByUri(&_Api.CallOpts, relay_uri)
}

// GetRelays is a free data retrieval call binding the contract method 0x04e469ea.
//
// Solidity: function get_relays() view returns((string,string,bool,string)[])
func (_Api *ApiCaller) GetRelays(opts *bind.CallOpts) ([]Struct0, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_relays")
	if err != nil {
		return *new([]Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new([]Struct0)).(*[]Struct0)

	return out0, err
}

// GetRelays is a free data retrieval call binding the contract method 0x04e469ea.
//
// Solidity: function get_relays() view returns((string,string,bool,string)[])
func (_Api *ApiSession) GetRelays() ([]Struct0, error) {
	return _Api.Contract.GetRelays(&_Api.CallOpts)
}

// GetRelays is a free data retrieval call binding the contract method 0x04e469ea.
//
// Solidity: function get_relays() view returns((string,string,bool,string)[])
func (_Api *ApiCallerSession) GetRelays() ([]Struct0, error) {
	return _Api.Contract.GetRelays(&_Api.CallOpts)
}

// GetRelaysAmount is a free data retrieval call binding the contract method 0x312c3165.
//
// Solidity: function get_relays_amount() view returns(uint256)
func (_Api *ApiCaller) GetRelaysAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Api.contract.Call(opts, &out, "get_relays_amount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetRelaysAmount is a free data retrieval call binding the contract method 0x312c3165.
//
// Solidity: function get_relays_amount() view returns(uint256)
func (_Api *ApiSession) GetRelaysAmount() (*big.Int, error) {
	return _Api.Contract.GetRelaysAmount(&_Api.CallOpts)
}

// GetRelaysAmount is a free data retrieval call binding the contract method 0x312c3165.
//
// Solidity: function get_relays_amount() view returns(uint256)
func (_Api *ApiCallerSession) GetRelaysAmount() (*big.Int, error) {
	return _Api.Contract.GetRelaysAmount(&_Api.CallOpts)
}

// AddRelay is a paid mutator transaction binding the contract method 0x2e21ecef.
//
// Solidity: function add_relay(string uri, string operator, bool is_mandatory, string description) returns()
func (_Api *ApiTransactor) AddRelay(opts *bind.TransactOpts, uri string, operator string, is_mandatory bool, description string) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "add_relay", uri, operator, is_mandatory, description)
}

// AddRelay is a paid mutator transaction binding the contract method 0x2e21ecef.
//
// Solidity: function add_relay(string uri, string operator, bool is_mandatory, string description) returns()
func (_Api *ApiSession) AddRelay(uri string, operator string, is_mandatory bool, description string) (*types.Transaction, error) {
	return _Api.Contract.AddRelay(&_Api.TransactOpts, uri, operator, is_mandatory, description)
}

// AddRelay is a paid mutator transaction binding the contract method 0x2e21ecef.
//
// Solidity: function add_relay(string uri, string operator, bool is_mandatory, string description) returns()
func (_Api *ApiTransactorSession) AddRelay(uri string, operator string, is_mandatory bool, description string) (*types.Transaction, error) {
	return _Api.Contract.AddRelay(&_Api.TransactOpts, uri, operator, is_mandatory, description)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0x253c8bd4.
//
// Solidity: function change_owner(address owner) returns()
func (_Api *ApiTransactor) ChangeOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "change_owner", owner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0x253c8bd4.
//
// Solidity: function change_owner(address owner) returns()
func (_Api *ApiSession) ChangeOwner(owner common.Address) (*types.Transaction, error) {
	return _Api.Contract.ChangeOwner(&_Api.TransactOpts, owner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0x253c8bd4.
//
// Solidity: function change_owner(address owner) returns()
func (_Api *ApiTransactorSession) ChangeOwner(owner common.Address) (*types.Transaction, error) {
	return _Api.Contract.ChangeOwner(&_Api.TransactOpts, owner)
}

// DismissManager is a paid mutator transaction binding the contract method 0x417a02b4.
//
// Solidity: function dismiss_manager() returns()
func (_Api *ApiTransactor) DismissManager(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "dismiss_manager")
}

// DismissManager is a paid mutator transaction binding the contract method 0x417a02b4.
//
// Solidity: function dismiss_manager() returns()
func (_Api *ApiSession) DismissManager() (*types.Transaction, error) {
	return _Api.Contract.DismissManager(&_Api.TransactOpts)
}

// DismissManager is a paid mutator transaction binding the contract method 0x417a02b4.
//
// Solidity: function dismiss_manager() returns()
func (_Api *ApiTransactorSession) DismissManager() (*types.Transaction, error) {
	return _Api.Contract.DismissManager(&_Api.TransactOpts)
}

// RecoverErc20 is a paid mutator transaction binding the contract method 0xedd885b4.
//
// Solidity: function recover_erc20(address token, uint256 amount, address recipient) returns()
func (_Api *ApiTransactor) RecoverErc20(opts *bind.TransactOpts, token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "recover_erc20", token, amount, recipient)
}

// RecoverErc20 is a paid mutator transaction binding the contract method 0xedd885b4.
//
// Solidity: function recover_erc20(address token, uint256 amount, address recipient) returns()
func (_Api *ApiSession) RecoverErc20(token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Api.Contract.RecoverErc20(&_Api.TransactOpts, token, amount, recipient)
}

// RecoverErc20 is a paid mutator transaction binding the contract method 0xedd885b4.
//
// Solidity: function recover_erc20(address token, uint256 amount, address recipient) returns()
func (_Api *ApiTransactorSession) RecoverErc20(token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Api.Contract.RecoverErc20(&_Api.TransactOpts, token, amount, recipient)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xf5a70a80.
//
// Solidity: function remove_relay(string uri) returns()
func (_Api *ApiTransactor) RemoveRelay(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "remove_relay", uri)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xf5a70a80.
//
// Solidity: function remove_relay(string uri) returns()
func (_Api *ApiSession) RemoveRelay(uri string) (*types.Transaction, error) {
	return _Api.Contract.RemoveRelay(&_Api.TransactOpts, uri)
}

// RemoveRelay is a paid mutator transaction binding the contract method 0xf5a70a80.
//
// Solidity: function remove_relay(string uri) returns()
func (_Api *ApiTransactorSession) RemoveRelay(uri string) (*types.Transaction, error) {
	return _Api.Contract.RemoveRelay(&_Api.TransactOpts, uri)
}

// SetManager is a paid mutator transaction binding the contract method 0x9aece83e.
//
// Solidity: function set_manager(address manager) returns()
func (_Api *ApiTransactor) SetManager(opts *bind.TransactOpts, manager common.Address) (*types.Transaction, error) {
	return _Api.contract.Transact(opts, "set_manager", manager)
}

// SetManager is a paid mutator transaction binding the contract method 0x9aece83e.
//
// Solidity: function set_manager(address manager) returns()
func (_Api *ApiSession) SetManager(manager common.Address) (*types.Transaction, error) {
	return _Api.Contract.SetManager(&_Api.TransactOpts, manager)
}

// SetManager is a paid mutator transaction binding the contract method 0x9aece83e.
//
// Solidity: function set_manager(address manager) returns()
func (_Api *ApiTransactorSession) SetManager(manager common.Address) (*types.Transaction, error) {
	return _Api.Contract.SetManager(&_Api.TransactOpts, manager)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_Api *ApiTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Api.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_Api *ApiSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Api.Contract.Fallback(&_Api.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_Api *ApiTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Api.Contract.Fallback(&_Api.TransactOpts, calldata)
}

// ApiAllowedListUpdatedIterator is returned from FilterAllowedListUpdated and is used to iterate over the raw logs and unpacked data for AllowedListUpdated events raised by the Api contract.
type ApiAllowedListUpdatedIterator struct {
	Event *ApiAllowedListUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiAllowedListUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiAllowedListUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiAllowedListUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiAllowedListUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiAllowedListUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiAllowedListUpdated represents a AllowedListUpdated event raised by the Api contract.
type ApiAllowedListUpdated struct {
	AllowedListVersion *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterAllowedListUpdated is a free log retrieval operation binding the contract event 0x49f5627aa055ec3fcd474f99c8b7799b798c04af7b9f215305512c867e5a1839.
//
// Solidity: event AllowedListUpdated(uint256 indexed allowed_list_version)
func (_Api *ApiFilterer) FilterAllowedListUpdated(opts *bind.FilterOpts, allowed_list_version []*big.Int) (*ApiAllowedListUpdatedIterator, error) {
	var allowed_list_versionRule []interface{}
	for _, allowed_list_versionItem := range allowed_list_version {
		allowed_list_versionRule = append(allowed_list_versionRule, allowed_list_versionItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "AllowedListUpdated", allowed_list_versionRule)
	if err != nil {
		return nil, err
	}
	return &ApiAllowedListUpdatedIterator{contract: _Api.contract, event: "AllowedListUpdated", logs: logs, sub: sub}, nil
}

// WatchAllowedListUpdated is a free log subscription operation binding the contract event 0x49f5627aa055ec3fcd474f99c8b7799b798c04af7b9f215305512c867e5a1839.
//
// Solidity: event AllowedListUpdated(uint256 indexed allowed_list_version)
func (_Api *ApiFilterer) WatchAllowedListUpdated(opts *bind.WatchOpts, sink chan<- *ApiAllowedListUpdated, allowed_list_version []*big.Int) (event.Subscription, error) {
	var allowed_list_versionRule []interface{}
	for _, allowed_list_versionItem := range allowed_list_version {
		allowed_list_versionRule = append(allowed_list_versionRule, allowed_list_versionItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "AllowedListUpdated", allowed_list_versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiAllowedListUpdated)
				if err := _Api.contract.UnpackLog(event, "AllowedListUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAllowedListUpdated is a log parse operation binding the contract event 0x49f5627aa055ec3fcd474f99c8b7799b798c04af7b9f215305512c867e5a1839.
//
// Solidity: event AllowedListUpdated(uint256 indexed allowed_list_version)
func (_Api *ApiFilterer) ParseAllowedListUpdated(log types.Log) (*ApiAllowedListUpdated, error) {
	event := new(ApiAllowedListUpdated)
	if err := _Api.contract.UnpackLog(event, "AllowedListUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApiERC20RecoveredIterator is returned from FilterERC20Recovered and is used to iterate over the raw logs and unpacked data for ERC20Recovered events raised by the Api contract.
type ApiERC20RecoveredIterator struct {
	Event *ApiERC20Recovered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiERC20RecoveredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiERC20Recovered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiERC20Recovered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiERC20RecoveredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiERC20RecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiERC20Recovered represents a ERC20Recovered event raised by the Api contract.
type ApiERC20Recovered struct {
	Token     common.Address
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterERC20Recovered is a free log retrieval operation binding the contract event 0x8619312ed4eff1cf9f0116e6db2f49d9570a86f0350d1c5ad1bd0f7b0cf9e132.
//
// Solidity: event ERC20Recovered(address indexed token, uint256 amount, address indexed recipient)
func (_Api *ApiFilterer) FilterERC20Recovered(opts *bind.FilterOpts, token []common.Address, recipient []common.Address) (*ApiERC20RecoveredIterator, error) {
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "ERC20Recovered", tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &ApiERC20RecoveredIterator{contract: _Api.contract, event: "ERC20Recovered", logs: logs, sub: sub}, nil
}

// WatchERC20Recovered is a free log subscription operation binding the contract event 0x8619312ed4eff1cf9f0116e6db2f49d9570a86f0350d1c5ad1bd0f7b0cf9e132.
//
// Solidity: event ERC20Recovered(address indexed token, uint256 amount, address indexed recipient)
func (_Api *ApiFilterer) WatchERC20Recovered(opts *bind.WatchOpts, sink chan<- *ApiERC20Recovered, token []common.Address, recipient []common.Address) (event.Subscription, error) {
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "ERC20Recovered", tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiERC20Recovered)
				if err := _Api.contract.UnpackLog(event, "ERC20Recovered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseERC20Recovered is a log parse operation binding the contract event 0x8619312ed4eff1cf9f0116e6db2f49d9570a86f0350d1c5ad1bd0f7b0cf9e132.
//
// Solidity: event ERC20Recovered(address indexed token, uint256 amount, address indexed recipient)
func (_Api *ApiFilterer) ParseERC20Recovered(log types.Log) (*ApiERC20Recovered, error) {
	event := new(ApiERC20Recovered)
	if err := _Api.contract.UnpackLog(event, "ERC20Recovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApiManagerChangedIterator is returned from FilterManagerChanged and is used to iterate over the raw logs and unpacked data for ManagerChanged events raised by the Api contract.
type ApiManagerChangedIterator struct {
	Event *ApiManagerChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiManagerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiManagerChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiManagerChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiManagerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiManagerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiManagerChanged represents a ManagerChanged event raised by the Api contract.
type ApiManagerChanged struct {
	NewManager common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterManagerChanged is a free log retrieval operation binding the contract event 0x198db6e425fb8aafd1823c6ca50be2d51e5764571a5ae0f0f21c6812e45def0b.
//
// Solidity: event ManagerChanged(address indexed new_manager)
func (_Api *ApiFilterer) FilterManagerChanged(opts *bind.FilterOpts, new_manager []common.Address) (*ApiManagerChangedIterator, error) {
	var new_managerRule []interface{}
	for _, new_managerItem := range new_manager {
		new_managerRule = append(new_managerRule, new_managerItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "ManagerChanged", new_managerRule)
	if err != nil {
		return nil, err
	}
	return &ApiManagerChangedIterator{contract: _Api.contract, event: "ManagerChanged", logs: logs, sub: sub}, nil
}

// WatchManagerChanged is a free log subscription operation binding the contract event 0x198db6e425fb8aafd1823c6ca50be2d51e5764571a5ae0f0f21c6812e45def0b.
//
// Solidity: event ManagerChanged(address indexed new_manager)
func (_Api *ApiFilterer) WatchManagerChanged(opts *bind.WatchOpts, sink chan<- *ApiManagerChanged, new_manager []common.Address) (event.Subscription, error) {
	var new_managerRule []interface{}
	for _, new_managerItem := range new_manager {
		new_managerRule = append(new_managerRule, new_managerItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "ManagerChanged", new_managerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiManagerChanged)
				if err := _Api.contract.UnpackLog(event, "ManagerChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseManagerChanged is a log parse operation binding the contract event 0x198db6e425fb8aafd1823c6ca50be2d51e5764571a5ae0f0f21c6812e45def0b.
//
// Solidity: event ManagerChanged(address indexed new_manager)
func (_Api *ApiFilterer) ParseManagerChanged(log types.Log) (*ApiManagerChanged, error) {
	event := new(ApiManagerChanged)
	if err := _Api.contract.UnpackLog(event, "ManagerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApiOwnerChangedIterator is returned from FilterOwnerChanged and is used to iterate over the raw logs and unpacked data for OwnerChanged events raised by the Api contract.
type ApiOwnerChangedIterator struct {
	Event *ApiOwnerChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiOwnerChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiOwnerChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiOwnerChanged represents a OwnerChanged event raised by the Api contract.
type ApiOwnerChanged struct {
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerChanged is a free log retrieval operation binding the contract event 0xa2ea9883a321a3e97b8266c2b078bfeec6d50c711ed71f874a90d500ae2eaf36.
//
// Solidity: event OwnerChanged(address indexed new_owner)
func (_Api *ApiFilterer) FilterOwnerChanged(opts *bind.FilterOpts, new_owner []common.Address) (*ApiOwnerChangedIterator, error) {
	var new_ownerRule []interface{}
	for _, new_ownerItem := range new_owner {
		new_ownerRule = append(new_ownerRule, new_ownerItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "OwnerChanged", new_ownerRule)
	if err != nil {
		return nil, err
	}
	return &ApiOwnerChangedIterator{contract: _Api.contract, event: "OwnerChanged", logs: logs, sub: sub}, nil
}

// WatchOwnerChanged is a free log subscription operation binding the contract event 0xa2ea9883a321a3e97b8266c2b078bfeec6d50c711ed71f874a90d500ae2eaf36.
//
// Solidity: event OwnerChanged(address indexed new_owner)
func (_Api *ApiFilterer) WatchOwnerChanged(opts *bind.WatchOpts, sink chan<- *ApiOwnerChanged, new_owner []common.Address) (event.Subscription, error) {
	var new_ownerRule []interface{}
	for _, new_ownerItem := range new_owner {
		new_ownerRule = append(new_ownerRule, new_ownerItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "OwnerChanged", new_ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiOwnerChanged)
				if err := _Api.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerChanged is a log parse operation binding the contract event 0xa2ea9883a321a3e97b8266c2b078bfeec6d50c711ed71f874a90d500ae2eaf36.
//
// Solidity: event OwnerChanged(address indexed new_owner)
func (_Api *ApiFilterer) ParseOwnerChanged(log types.Log) (*ApiOwnerChanged, error) {
	event := new(ApiOwnerChanged)
	if err := _Api.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApiRelayAddedIterator is returned from FilterRelayAdded and is used to iterate over the raw logs and unpacked data for RelayAdded events raised by the Api contract.
type ApiRelayAddedIterator struct {
	Event *ApiRelayAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiRelayAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiRelayAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiRelayAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiRelayAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiRelayAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiRelayAdded represents a RelayAdded event raised by the Api contract.
type ApiRelayAdded struct {
	UriHash common.Hash
	Relay   Struct0
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRelayAdded is a free log retrieval operation binding the contract event 0xeee5faa84d45af657ab405cdbf2c6a8a6d466e83fa694a358fee5ff84431d0bf.
//
// Solidity: event RelayAdded(string indexed uri_hash, (string,string,bool,string) relay)
func (_Api *ApiFilterer) FilterRelayAdded(opts *bind.FilterOpts, uri_hash []string) (*ApiRelayAddedIterator, error) {
	var uri_hashRule []interface{}
	for _, uri_hashItem := range uri_hash {
		uri_hashRule = append(uri_hashRule, uri_hashItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "RelayAdded", uri_hashRule)
	if err != nil {
		return nil, err
	}
	return &ApiRelayAddedIterator{contract: _Api.contract, event: "RelayAdded", logs: logs, sub: sub}, nil
}

// WatchRelayAdded is a free log subscription operation binding the contract event 0xeee5faa84d45af657ab405cdbf2c6a8a6d466e83fa694a358fee5ff84431d0bf.
//
// Solidity: event RelayAdded(string indexed uri_hash, (string,string,bool,string) relay)
func (_Api *ApiFilterer) WatchRelayAdded(opts *bind.WatchOpts, sink chan<- *ApiRelayAdded, uri_hash []string) (event.Subscription, error) {
	var uri_hashRule []interface{}
	for _, uri_hashItem := range uri_hash {
		uri_hashRule = append(uri_hashRule, uri_hashItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "RelayAdded", uri_hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiRelayAdded)
				if err := _Api.contract.UnpackLog(event, "RelayAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayAdded is a log parse operation binding the contract event 0xeee5faa84d45af657ab405cdbf2c6a8a6d466e83fa694a358fee5ff84431d0bf.
//
// Solidity: event RelayAdded(string indexed uri_hash, (string,string,bool,string) relay)
func (_Api *ApiFilterer) ParseRelayAdded(log types.Log) (*ApiRelayAdded, error) {
	event := new(ApiRelayAdded)
	if err := _Api.contract.UnpackLog(event, "RelayAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApiRelayRemovedIterator is returned from FilterRelayRemoved and is used to iterate over the raw logs and unpacked data for RelayRemoved events raised by the Api contract.
type ApiRelayRemovedIterator struct {
	Event *ApiRelayRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ApiRelayRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApiRelayRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ApiRelayRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ApiRelayRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApiRelayRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApiRelayRemoved represents a RelayRemoved event raised by the Api contract.
type ApiRelayRemoved struct {
	UriHash common.Hash
	Uri     string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRelayRemoved is a free log retrieval operation binding the contract event 0xef756634af7ee7210f786ec0f91930afa63fda84d9ff6493ae681c332055dadb.
//
// Solidity: event RelayRemoved(string indexed uri_hash, string uri)
func (_Api *ApiFilterer) FilterRelayRemoved(opts *bind.FilterOpts, uri_hash []string) (*ApiRelayRemovedIterator, error) {
	var uri_hashRule []interface{}
	for _, uri_hashItem := range uri_hash {
		uri_hashRule = append(uri_hashRule, uri_hashItem)
	}

	logs, sub, err := _Api.contract.FilterLogs(opts, "RelayRemoved", uri_hashRule)
	if err != nil {
		return nil, err
	}
	return &ApiRelayRemovedIterator{contract: _Api.contract, event: "RelayRemoved", logs: logs, sub: sub}, nil
}

// WatchRelayRemoved is a free log subscription operation binding the contract event 0xef756634af7ee7210f786ec0f91930afa63fda84d9ff6493ae681c332055dadb.
//
// Solidity: event RelayRemoved(string indexed uri_hash, string uri)
func (_Api *ApiFilterer) WatchRelayRemoved(opts *bind.WatchOpts, sink chan<- *ApiRelayRemoved, uri_hash []string) (event.Subscription, error) {
	var uri_hashRule []interface{}
	for _, uri_hashItem := range uri_hash {
		uri_hashRule = append(uri_hashRule, uri_hashItem)
	}

	logs, sub, err := _Api.contract.WatchLogs(opts, "RelayRemoved", uri_hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApiRelayRemoved)
				if err := _Api.contract.UnpackLog(event, "RelayRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayRemoved is a log parse operation binding the contract event 0xef756634af7ee7210f786ec0f91930afa63fda84d9ff6493ae681c332055dadb.
//
// Solidity: event RelayRemoved(string indexed uri_hash, string uri)
func (_Api *ApiFilterer) ParseRelayRemoved(log types.Log) (*ApiRelayRemoved, error) {
	event := new(ApiRelayRemoved)
	if err := _Api.contract.UnpackLog(event, "RelayRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

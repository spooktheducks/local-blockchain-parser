package cmds

import (
	// "encoding/json"
	// "encoding/hex"
	// "encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"

	"github.com/WikiLeaksFreedomForce/local-blockchain-parser/cmds/utils"
)

var cablegateTxs = map[string]bool{
	"5c593b7b71063a01f4128c98e36fb407b00a87454e67b39ad5f8820ebc1b2ad5": false,
	"221d900b5ac701028f9dfab7dfba326f608308386d45c05432e721b7c122cba7": false,
	"cadfd932a4be36b635b633bd2bf4a2b3de72041b13d8331360cc2b707c8ed27c": false,
	"b5995cca21585c1b83838d5c654c24b4ad25ef717563758ea38cb0f019f9fa50": false,
	"84af6f819851dff5622031274ff3265a532881eab4829a080d026eac3addbbdd": false,
	"b28a631c8df1380f5c28d59fb78d98af8f899d69f6d7824bbb51c3e1de135fc1": false,
	"985b117137eee4a028a86646a1d45bf1ee6d6246f363324c4e92c477deb7fb1c": false,
	"7f8e5d2b11638736a03488f5aba5c9cdde55740d0e568f8c8c3d46623b5f6855": false,
	"3b98c3020695f0b089ff20a83d06dfb697e6e79e33f5d44d05292b5e0d708643": false,
	"d9aa09de38f7217344a82bc36056f54c2db5690096ba0012660d72ac1ec1fe19": false,
	"5f0025cbb5015dfe84d317545cdc2e3ee9ece7bd1bbaa45aa1465b8fa094e609": false,
	"bbf2d6dcf82fdee6b8ed8371ccb96ee8112282b0b7cdaf9eb8ba5cbe2f0db4b4": false,
	"045727591618d2154239956332f307c754a52c0e7896856a05848ec7065db1c3": false,
	"c0f01afc54bbb4e481310743ecce49ad7e5d4e3ebe8dae5b906c771f35dbebd6": false,
	"145465584a76c7a12842693451f617b9178828962668ac469719324283756653": false,
	"493dde8af86c6d9a451524f0ac8e104e903fc93b7d1f4a758399cc6953b9ec33": false,
	"2c9e766020d9e93bea3a1d149313ab224d3c375ad9341594331fa9c48bce13b8": false,
	"296913b89c41b19811e4dff8f74779cd1ba28597c508775a431763d79e690004": false,
	"abe43334f768810918958953d38e5d5d04f5a5ce1b6b0bb209b358365a3a653a": false,
	"d5d9f7e16169b7f65f9bc7e6fed1635f61b1f33cdbec131894d48acc5df2a4b2": false,
	"281b42ecff42b0a4a684cb28a86336a14eb7ba0bf263f12e35d5f23cd3fdd705": false,
	"838c3d543f222c5cfa1d5395c0c376e0e586ffab5850ff7767ed1a2c522db4f5": false,
	"8b0371ebe3fa924d5c8d45fff5923c7514dcfcaeed426d1c792653aca5f41dd2": false,
	"4a2893f4ed4eb86af4eb5c16af3df39570fc40e1892bb2961cedd48cbd3b1a6c": false,
	"f480fa01de505925722bf49f9ac28d3f5e424fb74c3ef41a5b56aca88b99ec4b": false,
	"375827f4af67ad78e50677ec414f016c0d46699fb3ff1f59673c16a72c1346a5": false,
	"f24c7ef69eb38aecff62e0e3aa317818a3c42ce48838b3edb106b50a6f9d023e": false,
	"4f0df2e6c54af6dc90fab6b568db42a5885998c676cb5b91d8402014219ed992": false,
	"41999d6b50b31307cea87cb385655abe494396989e650ed6966ad1cb4c36dd0f": false,
	"d3eb7bc5a907b9bd3c5c59f85b21aa53377032f0a6fc54e5b211977d2e9aae2d": false,
	"f29ae9e17e7a0210090a674a14c93426f88569a351f10976415761a82a76f85e": false,
	"ce71b791ea2464e597161f971c8b18ea4291d16760c6c626af7bf98e3bb7dbc3": false,
	"ac196df899cb2c90840637a1a0bc13bbaa46740ac527adb43f7b4fe91686be07": false,
	"f1656e37b73c8c856e0288db5451fca24273f4a99f34dab43863578e2e834fb9": false,
	"0fb6b21b3c3aa80da16671a5b42509afc00e4e0772100117b2b1d2f3053a3489": false,
	"423774735e612bd52c91380a3d08fc4187ed06e0e75ebc2eaadf0bdb1f6838da": false,
	"942d24347fbd4775d6531ab947c557edbb64a155576895f75b0552ac9c87a61f": false,
	"daa62cf04e78e56acebe086565fc57a97f4abf11f433774da3920cf5b950fd07": false,
	"2343eac8264456f756720a6633025f622880c0c4ce1defb2511320ae6761eebe": false,
	"40ca3c944963f40b06769ade234e0207ccfa42921221d0244b25c39a1197fa21": false,
	"18781b3756f2cf11dc9a87e949fa14d51bf8f395482ae952adbc61106509370d": false,
	"63b6af81eba7cc84d0b33888ac044d455e80c86a807f5510d36ae290cb0d86cb": false,
	"953113160ed7c0cbaf7936751c02c01b90bce3dda05eb04512c34507ef947488": false,
	"7615f35301ad1998e3124c102ea6d97eeb3d70c43e33a9b9c708ca80cb306732": false,
	"45f4e9ef7ba876a0a5a67c32c6acd431c0541dd837ec651dbf2122c1e741a612": false,
	"257a87e108f2df3bd52fc5755f5fa8cbe6f3d0b7cc58c0dd6882520845d2ca96": false,
	"5a665ba2f149512e2d50f8b038506ba5025343448c35dcce5e7877847f892118": false,
	"bf156d43f3cf3af47c19873b4ff4202bc3f7984916dfb9100c379387708142ae": false,
	"4f47c8ca1c5d66d3635260c077dd2c3608e4c9cbae4956bf47c952671ee408be": false,
	"1dea410e76d9971c2c584e2f5a04c7620d960f88f0d60ae8e9a81aa4502fc0cd": false,
	"2e7d8069f3a2e907a412854930ac8c2216e38b89ba160ecc3bd55aee82054d4c": false,
	"6227cc7bc1d761bfb3f7984403d8a39465d694b9428173bb8ee9f64cb0fb0266": false,
	"cec970b5075a913168e97d5e4db94a4f97f37b11d365e2f2009f8fd96789461b": false,
	"79536c57b0311d836fc80e547787867d880f7eecc68d61c81a4d929648bd7a3e": false,
	"63038c313aad91adc520dc648f88a48d1e85a036226191982da2e9c1c8fbf003": false,
	"13e085a5a83ffc8fc1bc06d877b2f15b42eeb97b911bb310ec51a9524b4dc919": false,
	"6d77a2208f9701db714fbdb3a0f224390c8fcd285155e529345ebc96598b58e5": false,
	"5fbf5151d7ed614a8139272b2f822a80227016c3591357a8c553e8c1bc01ff1e": false,
	"e90a605ccf6a352213317ed5f26220a765e49b4f333a5b374d5c410d0f4467b9": false,
	"49f09eede6f7dc30dd9442c2f865338d25c2dc981246ca25288bf23343616040": false,
	"dfdf1772a7ecf8c634931b1fa8d38b88606a12a77583cb8ec8c19a4070739e6a": false,
	"d5da3645962491825de4e3325e0a57e966727dac0559d8ecb470b4449d77330e": false,
	"75b8f5f63c5328da5476b6ddc63196030b32b1d9d3216ccd1ca09505e701c158": false,
	"9fe44168ff583349588fffa48adfca0ed2d2a4521dfbd76af911416139aabef1": false,
	"791d912d047ae36fbaff58653e753b7d78ecfef155c01670f14fea3a0542d202": false,
	"83771b859679890037422f228b02f000f27a19309e78f9c0abacb8e32c501112": false,
	"5e1cb0c14ab74e876426b14b7a7982a7038103f543adf6eb0ffa06707c0bee5b": false,
	"649904ea8b9687cb42a56fd66d8e4bb20590420923fbb847bb6b5933b106df59": false,
	"85d8c9dd08d5f37e6328af52c55d4f1c0e94a298957bf5cfbac6266b2ef15c72": false,
	"3b69cb34c8390589ef988ee8f9d51bbddecf0d0a1452fc968754162c29baec37": false,
	"a3b5cf488097399e4f481e077d626a09b2ccc8e1070035139eef67995e96190c": false,
	"4dfe95337d637cfab7d6387c328e26764860ba3df60ef28f58ed4cdd95e2de60": false,
	"c70fe6decd534f81c79536b4099f738464b2948f17fc6e686178ef26792bbf97": false,
	"541d6f8c4c7138fcf109953ac23e2ee1015c4ddee5fe3c407a3173750aef30c8": false,
	"ce2f7d28eb3fcdd24c2dfb9bbe2fad54aeca3713e653f192ab65a888399cddee": false,
	"14b3657b7342dd3cea67298f2d8f79cfd64a621b1c2efb5cb6144afcd70a8217": false,
	"beb13851fd7bb8afc502ae2719226b1f0724c316adb3678b5a2c5cf91d18ff40": false,
	"707b56d010b37f5832916773a102e3cb2063756d5325f4fe484db88fd9ac4bb6": false,
	"6ade3a5abb7021920eb5fb3e5049742b3a7860345753e904ee3f0f65e5fe504d": false,
	"3cb41a98bbbcab83edd2ccf628cbcb04a0a5cfd11f5dcc8380a7ed3366becf1c": false,
	"5a1d6a2b2869d61a65cdb523c572f2b0be6f8deecac9011dd5a2d67eef198f41": false,
	"2b3eda4d2fa1561771eb88c7cbd063348d5bb7751d97a4bfe227fee8166a661b": false,
	"0cf5a6627e897c29ff10137da9603f56c9212011bcbcf57e8a84de57b4ea1816": false,
	"5615161847439419be5a06d0918f4f42f3d7354280ddba50f628d6bee9822b23": false,
	"bda98cbb9751d67d415a9a2d43029442c1edf4428d169531172e8306438c4a12": false,
	"7d9dcf0e419bda1f791343eaa8c0e8a29bed461b5d792acf94478d2c06a3ecc7": false,
	"bda734cab4b9ac1e7be8ee5fd7f6ac33af54acd4f193517aa88110eec3481ceb": false,
	"ca3979ebfea7b0956eeb32869cd4e0e354b5a6f613110c2562a665b773097600": false,
	"36205791eec697e8ca008841fa043e199cb17d813d7dd77789374ed58a086bda": false,
	"876df0845093910edfc59937b92f47e81e11a33c05a594545708c85e2b751f67": false,
	"5cd925dc8a139b822564d7b5f3b81f1430893c23262c151bb1e441cbe91f38cd": false,
	"91a9beae4416c3f8d3447ae79945848f28f87b415ce6448d99f5bf35c69c194c": false,
	"eb4ead788b851c29a3db656c3279a3a38fc5f333b937d81a25795dc87ecfb006": false,
	"ed7e1d6877d414daaf6cb2b7ce622ae60607620311febe1835e65e6868b5666f": false,
	"255651e88eded9bbd7dc1d7db208446b26dc19cf8a62d976d7f0e840837d48eb": false,
	"4937171876aa9779ceef1899a1a145b6d353ce8acde6822916f6df458720f0ee": false,
	"026bb25ab801ad6155c3dc0033741dcde85990c81689e96e64da899fc4700a1b": false,
	"2fec6c7177f157a9a2cbd7da9356fd9dadf0da29b3ce7ee5c0bbf4539c9d8002": false,
	"d6e54b152b4c4f0e850659d1ba07e02da37fbf115facc55e02467407885fd528": false,
	"0dc3ae4c9ca80d2fdeebcbe5163ff9a96bf8b6f3a4a27961dce7b20e168abec2": false,
	"6e9543698260d40e00064a8127928c5f274b7988caa8707d23853417b456799b": false,
	"8698820d00cd7b15702da2366a800fd4da53580d10d5388ead2b10f2cb62821a": false,
	"07e39d80d55614e78ea8eb69b4460ea91fbda08da10716dabad5ed3056b58fef": false,
	"9e2b7fa59d0cd0e7810bd3971c367421e5348a8b734c2da516a920ae792e75e6": false,
	"13892cdf33e84f3d536787e72b884882cc16412844cacce783d7a1afb0ea820c": false,
	"d79445765cb73a45912552dadf92dc692bf0fedfe512cce3a68444814bb00f51": false,
	"d7b5168b9f422bed77f9d116c6707dc76cf74d9618585cca154c2730ec4554c4": false,
	"03eaff1e6ceeb8c4c73f7946818c430f5bb87a9ff2e951f123e8445de1163a97": false,
	"9b319f18fac32824450c422f9eb79693db314cf118b6cc6231795d23c24731d8": false,
	"9a2ceab7c06abdebb522679ac4658ca211b4aea1f8e98b2446d06fad43b564ef": false,
	"f78ddd7a0e89cc5c11b284b2bb53c0fde59c231df5e774c421f87fb993342876": false,
	"c6dfc55675414c61cdaebae6629f333e0102ae8d5ce468a0679741330cff5e9a": false,
	"73a2b0574f1ead5240ce721ec485d0f2529b0c5d412fa48ea48446905d1c100f": false,
	"b9006e8d045c7b956223b5049d7f6414cbd800792c4102fc62bca05a0c65dfc6": false,
	"4595448695e4808ef84326dd2e622015b5f8cb22bdba521101f10fe6135a0170": false,
	"8f427ede9b3418f425740d0f934aa37da27adee59469de52a4765f25b668b71a": false,
	"b28c5053222b5dea42dff791890ed8907563da76cbab8e06e93ca22d250e747f": false,
	"bcbe10e738ac66c8a1734e01575312639a4b0e91cc06454d21d8741f5492f924": false,
	"91d5eb806e7d3d5b6fbe1a208df333c080a81bed292b3182068316ad347c55b2": false,
	"7600603412dfeef73d74a2666a5d4938ac574fe0ddba089a1d93ed99f8d3687b": false,
	"f0d03810b826305679086fb27cd87381fe3baf7e4fb65cf0317160e252f47a3b": false,
	"847209fc2e1912253783aca0bcf2f40ea5a936a3f7cbdaa31ec0a41817ea9f2c": false,
	"759591a0138c5b9e38677ba2a4bdff9a84c8e708054f5015a8916d62998cd9a7": false,
	"067bf093d9c947feda1d5f8294f1fa81be29235779eccd75388ab49f067500e6": false,
	"88a66cd01c966b69288cf5ea211d04b7c1480afcd1840b41e0717c9cc9260fdd": false,
	"690271f80625aef45a5f40369945d0f48681f712b47017bfef44b708963f1b88": false,
	"7478348820dbd554693afa52de07dddd09b2897cbd0d8874bce080f5d09a33f7": false,
	"cd43e05708b2c1ea5450d83ad4ee0b5418e4c3ec259c74dace6fd088b44cef8c": false,
	"c66bc1996c8b94a3811b365618fe1e9b1b56cacce993a3988a469fcc6d6dabfd": false,
	"2663cfa9cf4c03c609c593c3e91fede7029123dd42d25639d38a6cf50ab4cd44": false,
}

func FindSuspiciousTxs(startBlock, endBlock uint64, inDir, outDir string) error {
	outSubdir := filepath.Join(".", outDir, "blockdb")

	err := os.MkdirAll(outSubdir, 0777)
	if err != nil {
		return err
	}

	// start a goroutine to log errors
	chErr := make(chan error)
	go func() {
		for err := range chErr {
			fmt.Println("error:", err)
		}
	}()

	// start a goroutine for each .dat file being parsed
	chDones := []chan bool{}
	for i := int(startBlock); i < int(endBlock)+1; i++ {
		chDone := make(chan bool)
		go suspiciousTxsParseBlock(inDir, outSubdir, i, chErr, chDone)
		chDones = append(chDones, chDone)
	}

	// wait for all ops to complete
	for _, chDone := range chDones {
		<-chDone
	}

	// close error channel
	close(chErr)

	return nil
}

type suspiciousTx struct {
	Hash         string
	ParentHashes []string
}

func suspiciousTxsParseBlock(inDir string, outDir string, blockFileNum int, chErr chan error, chDone chan bool) {
	defer close(chDone)

	filename := fmt.Sprintf("blk%05d.dat", blockFileNum)
	fmt.Println("parsing block", filename)

	blocks, err := utils.LoadBlockFile(filepath.Join(inDir, filename))
	if err != nil {
		chErr <- err
		return
	}

	suspiciousTxs := map[string]suspiciousTx{}

	for _, bl := range blocks {
		// blockHash := bl.Hash().String()

		for _, tx := range bl.Transactions() {
			txHash := tx.Hash().String()

			if isSuspiciousTx(tx) {
				parentHashes := []string{}
				for _, txin := range tx.MsgTx().TxIn {
					parentHash := txin.PreviousOutPoint.Hash.String()
					parentHashes = append(parentHashes, parentHash)

					// x, err := hex.DecodeString(string(txin.SignatureScript[1:]))
					// if err != nil {
					//  fmt.Println("error:", err.Error())
					// }
					// fmt.Println("("+fmt.Sprintf("%v", len(txin.SignatureScript))+")", x)
					// fmt.Println(string(txin.SignatureScript))
					// ss := txin.SignatureScript
					// x, n := binary.Varint(ss[4:])
					// fmt.Printf("%v (len %v)\n", x, n)
					if txin.SignatureScript[0] == txscript.OP_PUSHDATA1 {
						fmt.Printf("OP_PUSHDATA1 (%v)\n", int(txin.SignatureScript[0]))
					} else if txin.SignatureScript[0] == txscript.OP_PUSHDATA2 {
						fmt.Println("OP_PUSHDATA2")
					} else if txin.SignatureScript[0] == txscript.OP_PUSHDATA4 {
						fmt.Println("OP_PUSHDATA4")
					}
				}

				suspiciousTxs[txHash] = suspiciousTx{
					Hash:         txHash,
					ParentHashes: parentHashes,
				}

			}
		}
	}

	// for txHash, tx := range suspiciousTxs {
	//  // if _, exists := suspiciousTxs[tx.ParentHash]; !exists {
	//  //  fmt.Println("found head ~>", txHash)
	//  // }
	//  fmt.Println(txHash)
	// for _, ph := range tx.ParentHashes {
	//  fmt.Println("  -", ph)
	// }
	// }

	for txHash := range cablegateTxs {
		fmt.Println(txHash)
		for _, ph := range suspiciousTxs[txHash].ParentHashes {
			fmt.Println("  -", ph)
		}
	}

	fmt.Println(len(cablegateTxs))
	fmt.Println(len(suspiciousTxs))
}

func isSuspiciousTx(tx *btcutil.Tx) bool {
	numTinyValues := 0
	for _, txout := range tx.MsgTx().TxOut {
		if utils.SatoshisToBTCs(txout.Value) == 0.00000001 {
			numTinyValues++
		}
	}

	if numTinyValues == len(tx.MsgTx().TxOut)-1 {
		return true
	}
	return false
}

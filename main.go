/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/ontio/avocado-transfer/common"
	sdk "github.com/ontio/ontology-go-sdk"
	ocommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

type Result struct {
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

func main() {
	err := common.DefConfig.Init("./config.json")
	if err != nil {
		fmt.Println("DefConfig.Init error:", err)
		return
	}

	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(common.DefConfig.JsonRpcAddress)
	user, ok := common.GetAccountByPassword(ontSdk, common.DefConfig.WalletFile)
	if !ok {
		fmt.Println("common.GetAccountByPassword error")
		return
	}

	var data []*Result
	var sts []*ont.State
	fi, err := os.Open(common.DefConfig.DataFile)
	if err != nil {
		fmt.Println("Error os.Open: ", err)
		return
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var sum uint64 = 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		result := new(Result)
		err := json.Unmarshal([]byte(a), result)
		if err != nil {
			fmt.Println("json.Unmarshal error")
			return
		}
		sum = sum + result.Value
		data = append(data, result)
	}

	f, err := os.Create("record.txt")
	if err != nil {
		fmt.Println("Error os.Create: ", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for _, record := range data {
		share := new(big.Int).SetUint64(record.Value)
		bonus := new(big.Int).SetUint64(common.DefConfig.Bonus)
		total := new(big.Int).SetUint64(sum)
		amount := new(big.Int).Div(new(big.Int).Mul(share, bonus), total)
		address, err := ocommon.AddressFromBase58(record.Address)
		if err != nil {
			fmt.Println("ocommon.AddressFromBase58 error:", record.Address)
			return
		}
		sts = append(sts, &ont.State{
			From:  user.Address,
			To:    address,
			Value: amount.Uint64(),
		})
		w.WriteString(address.ToBase58() + "\t" + amount.String())
		w.WriteString("\n")
	}
	w.Flush()

	n := (len(sts) + 499) / 500
	for i := 0; i < n; i++ {
		states := sts[:0]
		if i < (n - 1) {
			states = sts[i*500: (i+1)*500]
		} else {
			states = sts[i*500:]
		}
		fmt.Println(len(states))
		var tx *types.MutableTransaction
		if common.DefConfig.Asset == "ong" {
			tx, err = ontSdk.Native.Ong.NewMultiTransferTransaction(common.DefConfig.GasPrice, common.DefConfig.GasLimit, states)
			if err != nil {
				fmt.Println("ontSdk.Native.Ong.NewMultiTransferTransaction error :", err)
				return
			}
		} else if common.DefConfig.Asset == "ont" {
			tx, err = ontSdk.Native.Ont.NewMultiTransferTransaction(common.DefConfig.GasPrice, common.DefConfig.GasLimit, states)
			if err != nil {
				fmt.Println("ontSdk.Native.Ong.NewMultiTransferTransaction error :", err)
				return
			}
		} else if common.DefConfig.Asset == "oep4" {
			contract, err := ocommon.AddressFromHexString(common.DefConfig.ContractAddress)
			if err != nil {
				fmt.Println("ocommon.AddressFromHexString error :", err)
				return
			}
			var args []interface{}
			for _, st := range states {
				args = append(args, []interface{}{st.From, st.To, st.Value})
			}
			params := []interface{}{"transferMulti", args}
			tx, err = ontSdk.NeoVM.NewNeoVMInvokeTransaction(common.DefConfig.GasPrice, common.DefConfig.GasLimit,
				contract, params)
			if err != nil {
				fmt.Println("ontSdk.Native.Ong.NewMultiTransferTransaction error :", err)
				return
			}
		} else {
			fmt.Println("asset type not supported")
			return
		}

		err = ontSdk.SignToTransaction(tx, user)
		if err != nil {
			fmt.Println("ontSdk.SignToTransaction error :", err)
			return
		}
		txHash, err := ontSdk.SendTransaction(tx)
		if err != nil {
			fmt.Println("ontSdk.SendTransaction error :", err)
			return
		}
		fmt.Println("tx success, txHash is :", txHash.ToHexString())
	}
}

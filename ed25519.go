// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main


import (
    "time"
    "fmt"
    "os"
    "encoding/hex"
    "github.com/consensys/gnark/frontend"
    "github.com/consensys/gnark-crypto/ecc"
    "github.com/consensys/gnark/frontend/cs/r1cs"
    "github.com/consensys/gnark/backend/groth16"
    "gnark-ed25519/edwards_curve"
    "gnark-ed25519/sha512"
)

type Eddsa25519Circuit struct {
    M []frontend.Variable
    Pk []frontend.Variable
    Sig []frontend.Variable
}

func (circuit *Eddsa25519Circuit) Define(api frontend.API) error {
    c, err := edwards_curve.New[edwards_curve.Ed25519, edwards_curve.Ed25519Scalars](api)
    if err != nil {
        return err
    }
    edwards_curve.CheckValid(c, circuit.Sig, circuit.M, circuit.Pk)
    return nil
}


type Sha512Circuit struct {
    in []frontend.Variable `gnark:"in"`
    out []frontend.Variable `gnark:"out"`
}

func (circuit *Sha512Circuit) Define(api frontend.API) error {
    res := sha512.Sha512(api, circuit.in)
    if len(res) != 512 { panic("bad length") }
    for i := 0; i < 512; i++ {
        api.AssertIsEqual(res[i], circuit.out[i])
    }
    return nil
}



func main() {
    err := mainImpl()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

// func mainImpl() error {
//     in := bytesToBits([]byte("Succinct Labs"))
//     out := hexToBits("503ace098aa03f6feec1b5df0a38aee923f744a775508bc81f2b94ad139be297c2e8cd8c44af527b5d3f017a7fc929892c896604047e52e3f518924f52bff0dc")

//     myCircuit := Sha512Circuit{
//         in,
//         out,
//     }
//     fmt.Println(time.Now(), "compiling...")
//     r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &myCircuit)
//     if err != nil {
//         return err
//     }

//     assignment := &Sha512Circuit{
//         in,
//         out,
//     }
//     fmt.Println(time.Now(), "generating witness...")
//     witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
//     publicWitness, _ := witness.Public()
//     fmt.Println(time.Now(), "groth setup...")
//     pk, vk, err := groth16.Setup(r1cs)
//     fmt.Println(time.Now(), "groth prove...")
//     proof, err := groth16.Prove(r1cs, pk, witness)
//     fmt.Println(time.Now(), "groth verify...")
//     err = groth16.Verify(proof, vk, publicWitness)
//     if err != nil {
//         return err
//     }
//     fmt.Println(proof)
//     return nil
// }

func mainImpl() error {
    M := "53756363696e6374204c616273"
    Pk := "f7ec1c43f4de9d49556de87b86b26a98942cb078486fdb44de38b80864c39731"
    Sig := "35c323757c20640a294345c89c0bfcebe3d554fdb0c7b7a0bdb72222c531b1ec849fed99a053e0f5b02dd9a25bb6eb018885526d9f583cdbde0b1e9f6329da09"

    myCircuit := Eddsa25519Circuit{
        M: hexToBits(M),
        Pk: hexToBits(Pk),
        Sig: hexToBits(Sig),
    }
    fmt.Println(time.Now(), "compiling...")
    r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &myCircuit)
    if err != nil {
        return err
    }

    assignment := &Eddsa25519Circuit{
        M: hexToBits(M),
        Pk: hexToBits(Pk),
        Sig: hexToBits(Sig),
    }
    fmt.Println(time.Now(), "generating witness...")
    witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
    publicWitness, _ := witness.Public()
    fmt.Println(time.Now(), "groth setup...")
    pk, vk, err := groth16.Setup(r1cs)
    fmt.Println(time.Now(), "groth prove...")
    proof, err := groth16.Prove(r1cs, pk, witness)
    fmt.Println(time.Now(), "groth verify...")
    err = groth16.Verify(proof, vk, publicWitness)
    if err != nil {
        return err
    }
    fmt.Println(proof)
    return nil
}

func hexToBits(h string) []frontend.Variable {
    b, err := hex.DecodeString(h)
    if err != nil {
        panic(err)
    }
    result := make([]frontend.Variable, len(b) * 8)
    for i, v := range b {
        for j := 0; j < 8; j++ {
            if (v & (1 << j)) != 0 {
                result[i*8+j] = 1
            } else {
                result[i*8+j] = 0
            }
        }
    }
    return result
}

func bytesToBits(arr []byte) []frontend.Variable {
    result := make([]frontend.Variable, len(arr) * 8)
    for i, v := range arr {
        for j := 0; j < 8; j++ {
            if (v & (1 << (7-j))) != 0 {
                result[i*8+j] = 1
            } else {
                result[i*8+j] = 0
            }
        }
    }
    return result
}
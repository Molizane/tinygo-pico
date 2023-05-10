package main

import "machine"

func main() {
    i2c := machine.I2C0
    
    err := i2c.Configure(machine.I2CConfig{
        SCL: machine.P0_16,
        SDA: machine.P0_17,
    })
    
    if err != nil {
        println("could not configure I2C:", err)
        return
    }
    
    start_address := 8       // lower addresses are reserved to prevent conflicts with other protocols
    end_address := 119       // higher addresses unlock other modes, like 10-bit addressing
    
    for i := start_address; i <= end_address; i++ {
    w := []byte{0x75}
    r := make([]byte, 1)
    err = i2c.Tx(0x68, w, r)
    
    if err != nil {
        println("could not interact with I2C device:", err)
        return
    }
    
    println("WHO_AM_I:", r[0]) // prints "WHO_AM_I: 104"
}
}

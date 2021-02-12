package main

import (
	"fmt"

	"github.com/Dadido3/configdb"
	"github.com/spf13/viper"
)

var c, _ = configdb.New([]configdb.Storage{
	configdb.UseYAMLFile("/home/jose/go/src/github.com/josejuanmontiel/golang/viper/test.yaml"),
})

func Average(xs []float64) float64 {
	total := float64(0)
	for _, x := range xs {
		total += x
	}

	// fmt.Printf("Current Unix Time (Before): %v\n", time.Now().Unix())
	// time.Sleep(2 * time.Second)
	// fmt.Printf("Current Unix Time (After): %v\n", time.Now().Unix())

	return total / float64(len(xs))
}

func Viper() int {
	viper.SetConfigName("test") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./")   // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	nivel := viper.GetInt("level1.level2")
	viper.Set("level1.level2", 2)

	return nivel
}

func ConfigDb() int {
	var f int
	// Pass a pointer to any object you want to read from the internal tree at the given path ".box.width"
	err := c.Get(".level1.level2", &f)
	if err != nil {
		panic(err)
	}

	b := 2
	// Pass a boolean to be written at the path ".todo.WriteCode"
	err = c.Set(".level1.level2", b)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	Viper()
	ConfigDb()
}

package model

import (
	"flag"
	"fmt"
	"math"
)

var Conf *Configuration

var Models = []interface{}{
	&Student{},
}

type Configuration struct {
	//gorm Open method arg
	//currently only support mysql
	//user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	Server string
	//private key
	Key string

	//parameters for Water Masking algorithm
	Gamma, Nu, Xi uint
	//FiledNames consist of the columns typed float64
	//Wanting to more dynamic, I must to use reflect to get all the float columns
	//This is only a research method.Consequently, I use hard-writing directly.
	FiledNames    []string
}

func FlagConfInit() {
	confDatabase := flag.String("db", "stu", "")
	conUser := flag.String("user", "root", "database user")

	// the default password is written here for  testifying some experiments quickly.
	conPassword := flag.String("password", "toor", "the database password corresponding to the user")
	conIpAddress := flag.String("ip", "127.0.0.1", "the ip address of corresponding database")
	conPort := flag.String("port", "3306", "database port")

	privateKey := flag.String("key", "Oliver", "the private key")
	gamma := flag.Uint("gamma", 2, "gamma")
	nu := flag.Uint("nu", 5, "nu")
	xi := flag.Uint("xi", 3, "xi")
	//filedNames := flag.

	flag.Parse()

	Conf = &Configuration{}
	//default arg is :
	//root:toor@(127.0.0.1:3306)/stu
	Conf.Server = fmt.Sprintf("%s:%s@(%s:%s)/%s", *conUser, *conPassword, *conIpAddress, *conPort, *confDatabase)
	Conf.Key = *privateKey


	defaultFileds := []string{"Score1","Score2","Score3",
		"Score4","Score5",}
	var n = uint(math.Min(float64(*nu), float64(len(defaultFileds))))
	Conf.FiledNames = ([]string{"Score1","Score2","Score3",
		"Score4","Score5",})[:n]

	Conf.Gamma, Conf.Nu, Conf.Xi = *gamma, n, *xi
}

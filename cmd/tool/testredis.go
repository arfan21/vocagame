package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arfan21/vocagame/config"
	dbredis "github.com/arfan21/vocagame/pkg/db/redis"
	"github.com/urfave/cli/v2"
)

func TestRedis() *cli.Command {
	return &cli.Command{
		Name:  "testredis",
		Usage: "Test redis connection",
		Action: func(c *cli.Context) error {
			_, err := config.LoadConfig()
			if err != nil {
				return err
			}

			_, err = config.ParseConfig(config.GetViper())
			if err != nil {
				return err
			}

			dbRedis, err := dbredis.New()
			if err != nil {
				return err
			}

			err = dbredis.Set(context.Background(), dbRedis, "test-string", "test", 10*time.Second)
			if err != nil {
				return err
			}

			res1, err := dbredis.Get[string](context.Background(), dbRedis, "test-string")
			if err != nil {
				return err
			}

			fmt.Println("res1", res1)

			jsonByte, err := json.Marshal("test")
			if err != nil {
				return err
			}

			fmt.Println("string(jsonByte)", string(jsonByte))

			var res2 any
			err = json.Unmarshal(jsonByte, &res2)
			if err != nil {
				return err
			}

			fmt.Println("res2", res2)

			type TestStruct struct {
				Name string `json:"name"`
			}

			testStruct := TestStruct{
				Name: "test",
			}

			err = dbredis.Set(context.Background(), dbRedis, "test-struct", testStruct, 10*time.Second)
			if err != nil {
				return err
			}

			res3, err := dbredis.Get[TestStruct](context.Background(), dbRedis, "test-struct")
			if err != nil {
				return err
			}

			testStructByte, err := json.MarshalIndent(res3, " ", " ")
			if err != nil {
				return err
			}
			fmt.Println("res3", string(testStructByte))

			return nil
		},
	}
}

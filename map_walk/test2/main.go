package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func main() {
	content := `{
		"dev": {
			"__ini__": {
				"zone_id": 1,
				"sDataUrl": "http://172.26.144.19/dev/sdata.zip",
				"redisCount": "1",
				"redisAddrs": "172.26.144.20:6000",
				"redisDb": "0",
				"grpcWaitTime": "30",
				"ServerWeChatUrl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a502e1c3-d44a-48be-b606-88860c3cf215",
				"XiaoMeiWeChatUrl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=6afa3141-aa5e-402b-b8b6-5304b0ce7aed"
			},
			"appConfig": {
				"battle": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"service": "game",
						"services|game.type": "1",
						"services|game.path": "game"
					}
				},
				"battle_manager": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"service": "battle,win_battle",
						"services|battle.type": "4",
						"services|battle.path": "battle",
						"services|win_battle.type": "4",
						"services|win_battle.path": "win_battle"
					}
				},
				"chat": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"service": "game,gateway-battle",
						"services|game.type": "2",
						"services|game.path": "game",
						"services|gateway-battle.type": "1",
						"services|gateway-battle.path": "gateway-battle",
						"GlobalReloadable|check_name": "0"
					}
				},
				"game": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"service": "login,chat,game,battle_manager,cross_season,gateway-battle",
						"services|login.type": "2",
						"services|login.path": "login",
						"services|login.version": "/global",
						"services|chat.type": "1",
						"services|chat.path": "chat",
						"services|game.type": "2",
						"services|game.path": "game",
						"services|battle_manager.type": "1",
						"services|battle_manager.path": "battle_manager",
						"services|cross_season.type": "1",
						"services|cross_season.path": "cross_season",
						"services|cross_season.version": "/global",
						"services|gateway-battle.type": "1",
						"services|gateway-battle.path": "gateway-battle",
						"Reloadable|gm.enable": "1",
						"GlobalReloadable|check_name": "0"
					}
				},
				"gateway": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"HttpPort": "20000",
						"BattlePort": "11000",
						"service": "login,game,battle,win_battle",
						"services|game.type": "1",
						"services|game.path": "game",
						"services|login.type": "2",
						"services|login.path": "login",
						"services|login.version": "/global",
						"services|battle.type": "3",
						"services|battle.path": "battle",
						"services|win_battle.type": "3",
						"services|win_battle.path": "win_battle",
						"GlobalReloadable|check_name": "0"
					}
				},
				"login": {
					"__ini__": {
						"listenHost": "0.0.0.0",
						"listenPort": "10000",
						"service": "gateway",
						"services|gateway.type": "2",
						"services|gateway.path": "gateway",
						"Reloadable|role.list.protobuf": "1",
						"Reloadable|tmp.user.password": "4321",
						"GlobalReloadable|check_name": "0"
					}
				}
			},
			"config": {
				"battle": {
					"battle1": {
						"__ini__": {
							"id": "battle1"
						}
					}
				},
				"battle_manager": {
					"battle_manager1": {
						"__ini__": {
							"id": "battle_manager1"
						}
					}
				},
				"game": {
					"game1": {
						"__ini__": {
							"id": "game1"
						}
					}
				},
				"gateway": {
					"gateway1": {
						"__ini__": {
							"id": "gateway1"
						}
					}
				},
				"login": {
					"login1": {
						"__ini__": {
							"id": "login1"
						}
					}
				}
			},
			"fixConfig": {
				"battle": {
					"__ini__": {
						"dev:qa2.haidao:battle1": "1"
					}
				},
				"battle_manager": {
					"__ini__": {
						"dev:qa2.haidao:battle_manager1": "1"
					}
				}
			}
		},
		"services": {
			"battle": {
				"dev:qa2.haidao:battle1": "172.26.144.19:33458"
			},
			"battle_manager": {
				"dev:qa2.haidao:battle_manager1": "172.26.144.19:33458"
			},
			"chat": {
				"dev:qa2.haidao:chat1": "172.26.144.19:33458"
			},
			"cross_season": {
				"global:qa2.haidao:cross_season1": "172.26.144.19:33458"
			},
			"game": {
				"global:qa2.haidao:game1": "172.26.144.19:33458"
			},
			"gateway": {
				"dev:qa2.haidao:gateway1": "172.26.144.19:33458"
			},
			"gateway-battle": {
				"dev:qa2.haidao:gateway1": "172.26.144.19:33458"
			},
			"login": {
				"dev:qa2.haidao:login1": "172.26.144.19:33458"
			}
		}
	}`
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(content), &m); err != nil {
		panic(err)
	}
	out := map[string]map[string]string{}
	walk(reflect.ValueOf(m), "", out)
	//fmt.Println(m)
	fmt.Println(out)
}

func walk(v reflect.Value, path string, out map[string]map[string]string) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			walk(v.MapIndex(k), fmt.Sprintf("%s/%s", path, k), out)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		makeContent(fmt.Sprintf("%d", v.Int()), path, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		makeContent(fmt.Sprintf("%d", v.Uint()), path, out)
	case reflect.Float32, reflect.Float64:
		makeContent(fmt.Sprintf("%f", v.Float()), path, out)
	case reflect.String:
		makeContent(v.String(), path, out)
	default:
		makeContent(fmt.Sprintf("%v", v), path, out)
	}
}

func makeContent(v string, path string, out map[string]map[string]string) {
	if strings.Contains(path, "__ini__") {
		temp := strings.Split(path, "/__ini__/")
		path = temp[0]
		if _, ok := out[path]; !ok {
			out[path] = map[string]string{}
		}
		if strings.Contains(temp[1], "|") {
			vv := strings.Split(temp[1], "|")
			if _, ok := out[path][vv[0]]; !ok {
				out[path][vv[0]] = ""
			}
			out[path][vv[0]] = out[path][vv[0]] + fmt.Sprintf("%s=%s\n", vv[1], v)
		} else {
			if _, ok := out[path][""]; !ok {
				out[path][""] = ""
			}
			out[path][""] = out[path][""] + fmt.Sprintf("%s=%s\n", temp[1], v)
		}
	} else {
		if _, ok := out[path]; !ok {
			out[path] = map[string]string{}
		}
		out[path][""] = v
	}
}

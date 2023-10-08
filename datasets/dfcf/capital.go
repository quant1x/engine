package dfcf

import (
	"encoding/json"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/http"
	urlpkg "net/url"
	"strings"
	"time"
)

// 数据来源, 股本结构: https://emweb.securities.eastmoney.com/PC_HSF10/CapitalStockStructure/Index?type=web&code=SZ002849#
// 数据接口: https://emweb.securities.eastmoney.com/PC_HSF10/CapitalStockStructure/PageAjax?code=SZ002849
const (
	testCapitalSample = `{
    "xsjj": [
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "LIFT_DATE": "2023-11-30 00:00:00",
            "LIFT_NUM": 277835875,
            "LIFT_TYPE": "定向增发机构配售股份",
            "TOTAL_SHARES_RATIO": 53.17,
            "UNLIMITED_A_SHARES_RATIO": 117.03
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "LIFT_DATE": "2024-01-17 00:00:00",
            "LIFT_NUM": 3435000,
            "LIFT_TYPE": "股权激励限售股份",
            "TOTAL_SHARES_RATIO": 0.66,
            "UNLIMITED_A_SHARES_RATIO": 1.45
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "LIFT_DATE": "2025-01-17 00:00:00",
            "LIFT_NUM": 3426000,
            "LIFT_TYPE": "股权激励限售股份",
            "TOTAL_SHARES_RATIO": 0.66,
            "UNLIMITED_A_SHARES_RATIO": 1.44
        }
    ],
    "gbjg": [
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "NON_FREE_SHARES": null,
            "LIMITED_SHARES": 285106650,
            "UNLIMITED_SHARES": 237409149,
            "TOTAL_SHARES": 522515799,
            "LISTED_A_SHARES": 237409149,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "OTHER_FREE_SHARES": null,
            "NON_FREESHARES_RATIO": null,
            "LIMITED_SHARES_RATIO": 54.564216153012,
            "LISTED_SHARES_RATIO": 45.435783846988,
            "TOTAL_SHARES_RATIO": "100.00",
            "LISTED_A_RATIOPC": 100,
            "LISTED_B_RATIOPC": null,
            "LISTED_H_RATIOPC": null,
            "LISTED_OTHER_RATIOPC": null,
            "LISTED_SUM_RATIOPC": 100
        }
    ],
    "lngbbd": [
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-06-13 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LIMITED_SHARES": 285106650,
            "LIMITED_OTHARS": 285091875,
            "LIMITED_DOMESTIC_NATURAL": 56962512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 237409149,
            "LISTED_A_SHARES": 237409149,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 522515799,
            "LIMITED_A_SHARES": 285106650,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 14775,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "股权激励限售流通股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-06-09 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LIMITED_SHARES": 289800650,
            "LIMITED_OTHARS": 289785875,
            "LIMITED_DOMESTIC_NATURAL": 61656512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 232715149,
            "LISTED_A_SHARES": 232715149,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 522515799,
            "LIMITED_A_SHARES": 289800650,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 14775,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "高管股份变动"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-02-28 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LIMITED_SHARES": 309973400,
            "LIMITED_OTHARS": 289785875,
            "LIMITED_DOMESTIC_NATURAL": 61656512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 212542399,
            "LISTED_A_SHARES": 212542399,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 522515799,
            "LIMITED_A_SHARES": 309973400,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 20187525,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-12-31 00:00:00",
            "TOTAL_SHARES": 502123615,
            "LIMITED_SHARES": 309973400,
            "LIMITED_OTHARS": 289785875,
            "LIMITED_DOMESTIC_NATURAL": 61656512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 192150215,
            "LISTED_A_SHARES": 192150215,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 502123615,
            "LIMITED_A_SHARES": 309973400,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 20187525,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-12-27 00:00:00",
            "TOTAL_SHARES": 502108630,
            "LIMITED_SHARES": 309958625,
            "LIMITED_OTHARS": 289785875,
            "LIMITED_DOMESTIC_NATURAL": 61656512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 192150005,
            "LISTED_A_SHARES": 192150005,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 502108630,
            "LIMITED_A_SHARES": 309958625,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 20172750,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "回购"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-09-30 00:00:00",
            "TOTAL_SHARES": 502673630,
            "LIMITED_SHARES": 310523625,
            "LIMITED_OTHARS": 290350875,
            "LIMITED_DOMESTIC_NATURAL": 62221512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 192150005,
            "LISTED_A_SHARES": 192150005,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 502673630,
            "LIMITED_A_SHARES": 310523625,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 20172750,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市,高管股份变动"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-06-30 00:00:00",
            "TOTAL_SHARES": 502669471,
            "LIMITED_SHARES": 305480437,
            "LIMITED_OTHARS": 290350875,
            "LIMITED_DOMESTIC_NATURAL": 62221512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 197189034,
            "LISTED_A_SHARES": 197189034,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 502669471,
            "LIMITED_A_SHARES": 305480437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-03-31 00:00:00",
            "TOTAL_SHARES": 502653243,
            "LIMITED_SHARES": 305480437,
            "LIMITED_OTHARS": 290350875,
            "LIMITED_DOMESTIC_NATURAL": 62221512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 197172806,
            "LISTED_A_SHARES": 197172806,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 502653243,
            "LIMITED_A_SHARES": 305480437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-01-17 00:00:00",
            "TOTAL_SHARES": 500508037,
            "LIMITED_SHARES": 305480437,
            "LIMITED_OTHARS": 290350875,
            "LIMITED_DOMESTIC_NATURAL": 62221512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195027600,
            "LISTED_A_SHARES": 195027600,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 500508037,
            "LIMITED_A_SHARES": 305480437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "限制性股票"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2021-12-31 00:00:00",
            "TOTAL_SHARES": 487993037,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 277835875,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195027600,
            "LISTED_A_SHARES": 195027600,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487993037,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2021-09-30 00:00:00",
            "TOTAL_SHARES": 487992939,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 277835875,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195027502,
            "LISTED_A_SHARES": 195027502,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487992939,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2021-06-30 00:00:00",
            "TOTAL_SHARES": 487989984,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 277835875,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195024547,
            "LISTED_A_SHARES": 195024547,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487989984,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 228129363,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2021-03-31 00:00:00",
            "TOTAL_SHARES": 487989694,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 120554604,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": 157281271,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195024257,
            "LISTED_A_SHARES": 195024257,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487989694,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 70848092,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-12-31 00:00:00",
            "TOTAL_SHARES": 487987374,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 120554604,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": 157281271,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195021937,
            "LISTED_A_SHARES": 195021937,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487987374,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 70848092,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-11-27 00:00:00",
            "TOTAL_SHARES": 487984982,
            "LIMITED_SHARES": 292965437,
            "LIMITED_OTHARS": 120554604,
            "LIMITED_DOMESTIC_NATURAL": 49706512,
            "LIMITED_STATE_LEGAL": 157281271,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195019545,
            "LISTED_A_SHARES": 195019545,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 487984982,
            "LIMITED_A_SHARES": 292965437,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": 70848092,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "增发A股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-11-19 00:00:00",
            "TOTAL_SHARES": 210149107,
            "LIMITED_SHARES": 15129562,
            "LIMITED_OTHARS": null,
            "LIMITED_DOMESTIC_NATURAL": null,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195019545,
            "LISTED_A_SHARES": 195019545,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 210149107,
            "LIMITED_A_SHARES": 15129562,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": null,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-11-18 00:00:00",
            "TOTAL_SHARES": 210149045,
            "LIMITED_SHARES": 15129562,
            "LIMITED_OTHARS": null,
            "LIMITED_DOMESTIC_NATURAL": null,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195019483,
            "LISTED_A_SHARES": 195019483,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 210149045,
            "LIMITED_A_SHARES": 15129562,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": null,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-09-30 00:00:00",
            "TOTAL_SHARES": 210148781,
            "LIMITED_SHARES": 15129562,
            "LIMITED_OTHARS": null,
            "LIMITED_DOMESTIC_NATURAL": null,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195019219,
            "LISTED_A_SHARES": 195019219,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 210148781,
            "LIMITED_A_SHARES": 15129562,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": null,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-06-30 00:00:00",
            "TOTAL_SHARES": 210147613,
            "LIMITED_SHARES": 15129562,
            "LIMITED_OTHARS": null,
            "LIMITED_DOMESTIC_NATURAL": null,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 195018051,
            "LISTED_A_SHARES": 195018051,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 210147613,
            "LIMITED_A_SHARES": 15129562,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": null,
            "LOCK_SHARES": 15129562,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市,高管股份变动"
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2020-03-31 00:00:00",
            "TOTAL_SHARES": 210145189,
            "LIMITED_SHARES": 34777487,
            "LIMITED_OTHARS": null,
            "LIMITED_DOMESTIC_NATURAL": null,
            "LIMITED_STATE_LEGAL": null,
            "LIMITED_OVERSEAS_NOSTATE": null,
            "LIMITED_OVERSEAS_NATURAL": null,
            "UNLIMITED_SHARES": 175367702,
            "LISTED_A_SHARES": 175367702,
            "B_FREE_SHARE": null,
            "H_FREE_SHARE": null,
            "FREE_SHARES": 210145189,
            "LIMITED_A_SHARES": 34777487,
            "NON_FREE_SHARES": null,
            "LIMITED_B_SHARES": null,
            "OTHER_FREE_SHARES": null,
            "LIMITED_STATE_SHARES": null,
            "LIMITED_DOMESTIC_NOSTATE": null,
            "LOCK_SHARES": 34777487,
            "LIMITED_FOREIGN_SHARES": null,
            "LIMITED_H_SHARES": null,
            "SPONSOR_SHARES": null,
            "STATE_SPONSOR_SHARES": null,
            "SPONSOR_SOCIAL_SHARES": null,
            "RAISE_SHARES": null,
            "RAISE_STATE_SHARES": null,
            "RAISE_DOMESTIC_SHARES": null,
            "RAISE_OVERSEAS_SHARES": null,
            "CHANGE_REASON": "债转股上市"
        }
    ],
    "gbgc": [
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-06-13 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LISTED_A_SHARES": 237409149,
            "LIMITED_SHARES": 285106650
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-06-09 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LISTED_A_SHARES": 232715149,
            "LIMITED_SHARES": 289800650
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2023-02-28 00:00:00",
            "TOTAL_SHARES": 522515799,
            "LISTED_A_SHARES": 212542399,
            "LIMITED_SHARES": 309973400
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-12-31 00:00:00",
            "TOTAL_SHARES": 502123615,
            "LISTED_A_SHARES": 192150215,
            "LIMITED_SHARES": 309973400
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-12-27 00:00:00",
            "TOTAL_SHARES": 502108630,
            "LISTED_A_SHARES": 192150005,
            "LIMITED_SHARES": 309958625
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-09-30 00:00:00",
            "TOTAL_SHARES": 502673630,
            "LISTED_A_SHARES": 192150005,
            "LIMITED_SHARES": 310523625
        },
        {
            "SECUCODE": "002758.SZ",
            "SECURITY_CODE": "002758",
            "END_DATE": "2022-06-30 00:00:00",
            "TOTAL_SHARES": 502669471,
            "LISTED_A_SHARES": 197189034,
            "LIMITED_SHARES": 305480437
        }
    ]
}`
)

type rawCapital struct {
	Xsjj []struct {
		SECUCODE              string  `json:"SECUCODE"`
		SECURITYCODE          string  `json:"SECURITY_CODE"`
		LIFTDATE              string  `json:"LIFT_DATE"`
		LIFTNUM               int     `json:"LIFT_NUM"`
		LIFTTYPE              string  `json:"LIFT_TYPE"`
		TOTALSHARESRATIO      float64 `json:"TOTAL_SHARES_RATIO"`
		UNLIMITEDASHARESRATIO float64 `json:"UNLIMITED_A_SHARES_RATIO"`
	} `json:"xsjj"` // 限售解禁
	Gbjg []struct {
		SECUCODE           string  `json:"SECUCODE"`
		SECURITYCODE       string  `json:"SECURITY_CODE"`
		NONFREESHARES      int     `json:"NON_FREE_SHARES"`
		LIMITEDSHARES      int     `json:"LIMITED_SHARES"`
		UNLIMITEDSHARES    int     `json:"UNLIMITED_SHARES"`
		TOTALSHARES        int     `json:"TOTAL_SHARES"`
		LISTEDASHARES      int     `json:"LISTED_A_SHARES"`
		BFREESHARE         int     `json:"B_FREE_SHARE"`
		HFREESHARE         int     `json:"H_FREE_SHARE"`
		OTHERFREESHARES    int     `json:"OTHER_FREE_SHARES"`
		NONFREESHARESRATIO int     `json:"NON_FREESHARES_RATIO"`
		LIMITEDSHARESRATIO float64 `json:"LIMITED_SHARES_RATIO"`
		LISTEDSHARESRATIO  float64 `json:"LISTED_SHARES_RATIO"`
		TOTALSHARESRATIO   string  `json:"TOTAL_SHARES_RATIO"`
		LISTEDARATIOPC     float64 `json:"LISTED_A_RATIOPC"`
		LISTEDBRATIOPC     float64 `json:"LISTED_B_RATIOPC"`
		LISTEDHRATIOPC     float64 `json:"LISTED_H_RATIOPC"`
		LISTEDOTHERRATIOPC float64 `json:"LISTED_OTHER_RATIOPC"`
		LISTEDSUMRATIOPC   int     `json:"LISTED_SUM_RATIOPC"`
	} `json:"gbjg"` // 股本结构
	Lngbbd []struct {
		SECUCODE               string      `json:"SECUCODE"`
		SECURITYCODE           string      `json:"SECURITY_CODE"`
		ENDDATE                string      `json:"END_DATE"`
		TOTALSHARES            int         `json:"TOTAL_SHARES"`
		LIMITEDSHARES          int         `json:"LIMITED_SHARES"`
		LIMITEDOTHARS          int         `json:"LIMITED_OTHARS"`
		LIMITEDDOMESTICNATURAL int         `json:"LIMITED_DOMESTIC_NATURAL"`
		LIMITEDSTATELEGAL      int         `json:"LIMITED_STATE_LEGAL"`
		LIMITEDOVERSEASNOSTATE interface{} `json:"LIMITED_OVERSEAS_NOSTATE"`
		LIMITEDOVERSEASNATURAL interface{} `json:"LIMITED_OVERSEAS_NATURAL"`
		UNLIMITEDSHARES        int         `json:"UNLIMITED_SHARES"`
		LISTEDASHARES          int         `json:"LISTED_A_SHARES"`
		BFREESHARE             interface{} `json:"B_FREE_SHARE"`
		HFREESHARE             interface{} `json:"H_FREE_SHARE"`
		FREESHARES             int         `json:"FREE_SHARES"`
		LIMITEDASHARES         int         `json:"LIMITED_A_SHARES"`
		NONFREESHARES          interface{} `json:"NON_FREE_SHARES"`
		LIMITEDBSHARES         interface{} `json:"LIMITED_B_SHARES"`
		OTHERFREESHARES        interface{} `json:"OTHER_FREE_SHARES"`
		LIMITEDSTATESHARES     interface{} `json:"LIMITED_STATE_SHARES"`
		LIMITEDDOMESTICNOSTATE int         `json:"LIMITED_DOMESTIC_NOSTATE"`
		LOCKSHARES             int         `json:"LOCK_SHARES"`
		LIMITEDFOREIGNSHARES   interface{} `json:"LIMITED_FOREIGN_SHARES"`
		LIMITEDHSHARES         interface{} `json:"LIMITED_H_SHARES"`
		SPONSORSHARES          interface{} `json:"SPONSOR_SHARES"`
		STATESPONSORSHARES     interface{} `json:"STATE_SPONSOR_SHARES"`
		SPONSORSOCIALSHARES    interface{} `json:"SPONSOR_SOCIAL_SHARES"`
		RAISESHARES            interface{} `json:"RAISE_SHARES"`
		RAISESTATESHARES       interface{} `json:"RAISE_STATE_SHARES"`
		RAISEDOMESTICSHARES    interface{} `json:"RAISE_DOMESTIC_SHARES"`
		RAISEOVERSEASSHARES    interface{} `json:"RAISE_OVERSEAS_SHARES"`
		CHANGEREASON           string      `json:"CHANGE_REASON"`
	} `json:"lngbbd"` // 历年股本变动
	Gbgc []struct {
		SECUCODE      string `json:"SECUCODE"`
		SECURITYCODE  string `json:"SECURITY_CODE"`
		ENDDATE       string `json:"END_DATE"`
		TOTALSHARES   int    `json:"TOTAL_SHARES"`
		LISTEDASHARES int    `json:"LISTED_A_SHARES"`
		LIMITEDSHARES int    `json:"LIMITED_SHARES"`
	} `json:"gbgc"` // 股本构成
}

const (
	urlCapitalStockStructure = "https://emweb.securities.eastmoney.com/PC_HSF10/CapitalStockStructure/PageAjax"
)

type StockCapital struct {
	Code            string // 证券代码
	Date            string // 变动日期
	TotalShares     int    // 总股本
	UnlimitedShares int    // 已流通股本
	ListedAShares   int    // 已上市流通A股
	ChangeReson     string // 变动原因
	UpdateTime      string // 更新时间
}

// CapitalChange 获取股本变动记录
//
//	deprecated: 不推荐, 太慢
func CapitalChange(securityCode string) (list []StockCapital) {
	code := proto.CorrectSecurityCode(securityCode)
	params := urlpkg.Values{
		"code": {strings.ToUpper(code)},
	}

	url := urlCapitalStockStructure + "?" + params.Encode()
	data, lastModified, err := http.Request(url, http.GET)
	//fmt.Println(api.Bytes2String(data))
	if err != nil {
		return
	}
	var css rawCapital
	err = json.Unmarshal(data, &css)
	if err != nil {
		return
	}
	if lastModified.UnixMilli() > 0 {
		lastModified = time.Now()
	}
	updateTime := lastModified.Format(time.DateTime)
	for _, v := range css.Lngbbd {
		sc := StockCapital{
			Code:            securityCode,
			Date:            trading.FixTradeDate(v.ENDDATE),
			TotalShares:     v.TOTALSHARES,
			UnlimitedShares: v.UNLIMITEDSHARES,
			ListedAShares:   v.LISTEDASHARES,
			ChangeReson:     v.CHANGEREASON,
			UpdateTime:      updateTime,
		}
		list = append(list, sc)
	}
	api.SliceSort(list, func(a, b StockCapital) bool {
		return a.Date > b.Date
	})
	return
}

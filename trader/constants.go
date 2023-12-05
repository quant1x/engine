package trader

const (
	qmtPositionsPath     = "qmt"           // 持仓缓存路径
	qmtPositionsFilename = "positions.csv" // 持仓数据文件名
)

// 账号类型
const (
	FUTURE_ACCOUNT        = 1     // 期货
	SECURITY_ACCOUNT      = 2     // 股票
	CREDIT_ACCOUNT        = 3     // 信用
	FUTURE_OPTION_ACCOUNT = 5     // 期货期权
	STOCK_OPTION_ACCOUNT  = 6     // 股票期权
	HUGANGTONG_ACCOUNT    = 7     // 沪港通
	INCOME_SWAP_ACCOUNT   = 8     // 美股收益互换
	NEW3BOARD_ACCOUNT     = 10    // 全国股转账号
	SHENGANGTONG_ACCOUNT  = 11    // 深港通
	AT_OFFSITEBANKING     = 13    // 场外理财账户
	AT_OUTTER_FUTURE      = 1001  // 期货外盘
	AT_IB                 = 1002  // IB
	AT_NS_TRUSTBANK       = 15001 // 场外托管
	AT_INTERBANK          = 15002 // 银行间账号
	AT_BANK               = 15003 // 银行账号
	AT_OTC                = 15005 // 场外账号
)

// AccountStatus 账号状态
type AccountStatus int

const (
	ACCOUNT_STATUS_INVALID       = -1 // 无效
	ACCOUNT_STATUS_OK            = 0  // 正常
	ACCOUNT_STATUS_WAITING_LOGIN = 1  // 连接中
	ACCOUNT_STATUSING            = 2  // 登陆中
	ACCOUNT_STATUS_FAIL          = 3  // 失败
	ACCOUNT_STATUS_INITING       = 4  // 初始化中
	ACCOUNT_STATUS_CORRECTING    = 5  // 数据刷新校正中
	ACCOUNT_STATUS_CLOSED        = 6  // 收盘后
	ACCOUNT_STATUS_ASSIS_FAIL    = 7  // 穿透副链接断开
	ACCOUNT_STATUS_DISABLEBYSYS  = 8  // 系统停用(总线使用-密码错误超限)
	ACCOUNT_STATUS_DISABLEBYUSER = 9  // 用户停用(总线使用)
)

// 交易方向
const (
	STOCK_BUY                        = 23
	STOCK_SELL                       = 24
	CREDIT_BUY                       = 23 // 担保品买入
	CREDIT_SELL                      = 24 // 担保品卖出
	CREDIT_FIN_BUY                   = 27 // 融资买入
	CREDIT_SLO_SELL                  = 28 // 融券卖出
	CREDIT_BUY_SECU_REPAY            = 29 // 买券还券
	CREDIT_DIRECT_SECU_REPAY         = 30 // 直接还券
	CREDIT_SELL_SECU_REPAY           = 31 // 卖券还款
	CREDIT_DIRECT_CASH_REPAY         = 32 // 直接还款
	CREDIT_FIN_BUY_SPECIAL           = 40 // 专项融资买入
	CREDIT_SLO_SELL_SPECIAL          = 41 // 专项融券卖出
	CREDIT_BUY_SECU_REPAY_SPECIAL    = 42 // 专项买券还券
	CREDIT_DIRECT_SECU_REPAY_SPECIAL = 43 // 专项直接还券
	CREDIT_SELL_SECU_REPAY_SPECIAL   = 44 // 专项卖券还款
	CREDIT_DIRECT_CASH_REPAY_SPECIAL = 45 // 专项直接还款
)

// 报价类型
const (
	LATEST_PRICE                  = 5  // 最新价
	FIX_PRICE                     = 11 // 指定价/限价
	MARKET_SH_CONVERT_5_CANCEL    = 42 // 最优五档即时成交剩余撤销[上交所][股票]
	MARKET_SH_CONVERT_5_LIMIT     = 43 // 最优五档即时成交剩转限价[上交所][股票]
	MARKET_PEER_PRICE_FIRST       = 44 // 对手方最优价格委托[上交所[股票]][深交所[股票][期权]]
	MARKET_MINE_PRICE_FIRST       = 45 // 本方最优价格委托[上交所[股票]][深交所[股票][期权]]
	MARKET_SZ_INSTBUSI_RESTCANCEL = 46 // 即时成交剩余撤销委托[深交所][股票][期权]
	MARKET_SZ_CONVERT_5_CANCEL    = 47 // 最优五档即时成交剩余撤销[深交所][股票][期权]
	MARKET_SZ_FULL_OR_CANCEL      = 48 // 全额成交或撤销委托[深交所][股票][期权]
)

// OrderStatus 委托状态
type OrderStatus int

const (
	ORDER_UNREPORTED      = 48  // 未报
	ORDER_WAIT_REPORTING  = 49  // 待报
	ORDER_REPORTED        = 50  // 已报
	ORDER_REPORTED_CANCEL = 51  // 已报待撤
	ORDER_PARTSUCC_CANCEL = 52  // 部成待撤
	ORDER_PART_CANCEL     = 53  // 部撤
	ORDER_CANCELED        = 54  // 已撤
	ORDER_PART_SUCC       = 55  // 部成
	ORDER_SUCCEEDED       = 56  // 已成
	ORDER_JUNK            = 57  // 废单
	ORDER_UNKNOWN         = 255 // 未知
)

package trader

// AccountType 账号类型
type AccountType = int

const (
	FUTURE_ACCOUNT        AccountType = 1     // 期货
	SECURITY_ACCOUNT      AccountType = 2     // 股票
	CREDIT_ACCOUNT        AccountType = 3     // 信用
	FUTURE_OPTION_ACCOUNT AccountType = 5     // 期货期权
	STOCK_OPTION_ACCOUNT  AccountType = 6     // 股票期权
	HUGANGTONG_ACCOUNT    AccountType = 7     // 沪港通
	INCOME_SWAP_ACCOUNT   AccountType = 8     // 美股收益互换
	NEW3BOARD_ACCOUNT     AccountType = 10    // 全国股转账号
	SHENGANGTONG_ACCOUNT  AccountType = 11    // 深港通
	AT_OFFSITEBANKING     AccountType = 13    // 场外理财账户
	AT_OUTTER_FUTURE      AccountType = 1001  // 期货外盘
	AT_IB                 AccountType = 1002  // IB
	AT_NS_TRUSTBANK       AccountType = 15001 // 场外托管
	AT_INTERBANK          AccountType = 15002 // 银行间账号
	AT_BANK               AccountType = 15003 // 银行账号
	AT_OTC                AccountType = 15005 // 场外账号
)

// AccountStatus 账号状态
type AccountStatus = int

const (
	ACCOUNT_STATUS_INVALID       AccountStatus = -1 // 无效
	ACCOUNT_STATUS_OK            AccountStatus = 0  // 正常
	ACCOUNT_STATUS_WAITING_LOGIN AccountStatus = 1  // 连接中
	ACCOUNT_STATUSING            AccountStatus = 2  // 登陆中
	ACCOUNT_STATUS_FAIL          AccountStatus = 3  // 失败
	ACCOUNT_STATUS_INITING       AccountStatus = 4  // 初始化中
	ACCOUNT_STATUS_CORRECTING    AccountStatus = 5  // 数据刷新校正中
	ACCOUNT_STATUS_CLOSED        AccountStatus = 6  // 收盘后
	ACCOUNT_STATUS_ASSIS_FAIL    AccountStatus = 7  // 穿透副链接断开
	ACCOUNT_STATUS_DISABLEBYSYS  AccountStatus = 8  // 系统停用(总线使用-密码错误超限)
	ACCOUNT_STATUS_DISABLEBYUSER AccountStatus = 9  // 用户停用(总线使用)
)

// OrderType 订单类型
type OrderType = int

const (
	STOCK_BUY                        OrderType = 23
	STOCK_SELL                       OrderType = 24
	CREDIT_BUY                       OrderType = 23 // 担保品买入
	CREDIT_SELL                      OrderType = 24 // 担保品卖出
	CREDIT_FIN_BUY                   OrderType = 27 // 融资买入
	CREDIT_SLO_SELL                  OrderType = 28 // 融券卖出
	CREDIT_BUY_SECU_REPAY            OrderType = 29 // 买券还券
	CREDIT_DIRECT_SECU_REPAY         OrderType = 30 // 直接还券
	CREDIT_SELL_SECU_REPAY           OrderType = 31 // 卖券还款
	CREDIT_DIRECT_CASH_REPAY         OrderType = 32 // 直接还款
	CREDIT_FIN_BUY_SPECIAL           OrderType = 40 // 专项融资买入
	CREDIT_SLO_SELL_SPECIAL          OrderType = 41 // 专项融券卖出
	CREDIT_BUY_SECU_REPAY_SPECIAL    OrderType = 42 // 专项买券还券
	CREDIT_DIRECT_SECU_REPAY_SPECIAL OrderType = 43 // 专项直接还券
	CREDIT_SELL_SECU_REPAY_SPECIAL   OrderType = 44 // 专项卖券还款
	CREDIT_DIRECT_CASH_REPAY_SPECIAL OrderType = 45 // 专项直接还款
)

// PriceType 报价类型
type PriceType = int

const (
	LATEST_PRICE                  PriceType = 5  // 最新价
	FIX_PRICE                     PriceType = 11 // 指定价/限价
	MARKET_SH_CONVERT_5_CANCEL    PriceType = 42 // 最优五档即时成交剩余撤销[上交所][股票]
	MARKET_SH_CONVERT_5_LIMIT     PriceType = 43 // 最优五档即时成交剩转限价[上交所][股票]
	MARKET_PEER_PRICE_FIRST       PriceType = 44 // 对手方最优价格委托[上交所[股票]][深交所[股票][期权]]
	MARKET_MINE_PRICE_FIRST       PriceType = 45 // 本方最优价格委托[上交所[股票]][深交所[股票][期权]]
	MARKET_SZ_INSTBUSI_RESTCANCEL PriceType = 46 // 即时成交剩余撤销委托[深交所][股票][期权]
	MARKET_SZ_CONVERT_5_CANCEL    PriceType = 47 // 最优五档即时成交剩余撤销[深交所][股票][期权]
	MARKET_SZ_FULL_OR_CANCEL      PriceType = 48 // 全额成交或撤销委托[深交所][股票][期权]
)

// OrderStatus 委托状态
type OrderStatus = int

const (
	ORDER_UNREPORTED      OrderStatus = 48  // 未报
	ORDER_WAIT_REPORTING  OrderStatus = 49  // 待报
	ORDER_REPORTED        OrderStatus = 50  // 已报
	ORDER_REPORTED_CANCEL OrderStatus = 51  // 已报待撤
	ORDER_PARTSUCC_CANCEL OrderStatus = 52  // 部成待撤
	ORDER_PART_CANCEL     OrderStatus = 53  // 部撤
	ORDER_CANCELED        OrderStatus = 54  // 已撤
	ORDER_PART_SUCC       OrderStatus = 55  // 部成
	ORDER_SUCCEEDED       OrderStatus = 56  // 已成
	ORDER_JUNK            OrderStatus = 57  // 废单
	ORDER_UNKNOWN         OrderStatus = 255 // 未知
)

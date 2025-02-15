# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.9.0] - 2025-02-15
### Changed
- 更新依赖库gotdx版本号到1.23.0

## [1.8.46] - 2024-12-27
### Changed
- 更新依赖库gotdx版本号到1.22.23

## [1.8.45] - 2024-08-06
### Changed
- 更新依赖库ta-lib版本号到0.7.27
- update changelog
- update changelog

## [1.8.44] - 2024-08-06
### Changed
- 更新依赖库ta-lib版本号到0.7.26
- update changelog

## [1.8.43] - 2024-08-06
### Changed
- 更新依赖库版本
- update changelog

## [1.8.42] - 2024-08-06
### Changed
- 加载缓存增加TODO项,确认和解决内存泄漏的问题
- 更新依赖库版本
- update changelog

## [1.8.41] - 2024-07-31
### Changed
- 特征旋转适配器get接口增加锁机制
- 修订QMT的定义为miniQMT
- 修订QMT的定义为miniQMT
- 新增实时测试流程图
- 修订早盘抢筹时序图
- 修订尾盘先手时序图
- 新增龙虎榜接口和结构
- 修订东方财富泛型响应结构体
- 黑白名单命令字增加显示列表功能

## [1.8.40] - 2024-07-14
### Changed
- 更新ta-lib版本到0.7.23
- 调整实时K线重试部分代码
- update changelog

## [1.8.39] - 2024-07-05
### Changed
- git仓库忽略html文件
- 新增go-echarts html图表测试代码
- 调整私有变量名, 解除和go-echarts包冲突
- 调整html图表的title
- 修订除权除息中的成交量复权中的描述
- 删除废弃的代码
- 删除废弃的代码
- 新增加一处除权除息前一个交易日的跟踪点
- 更新exchange版本号到0.5.9
- 调整部分测试代码
- update changelog

## [1.8.38] - 2024-06-27
### Changed
- 更新依赖库版本
- 更新gotdx版本到1.22.20, 修复日线多次除权出息后,股价因为精度损失存在不准确的bug
- update changelog

## [1.8.37] - 2024-06-27
### Changed
- 更新gotdx版本到1.22.19
- 更新日K线, 在交易日9点15分之前, 屏蔽交易日的数据, 以保持除权的有效状态, 即只有1条上一个交易日的数据
- update changelog

## [1.8.36] - 2024-06-26
### Changed
- 修复除权bug,T+1除权,T日除权正确,T日前的数据存在重复除权的问题
- update changelog

## [1.8.35] - 2024-06-26
### Changed
- 更新ta-lib版本号到0.7.18
- update changelog

## [1.8.34] - 2024-06-26
### Changed
- 更新gotdx版本到1.22.18
- update changelog

## [1.8.33] - 2024-06-25
### Changed
- tick模式午间休市, tracker不退出
- update changelog

## [1.8.32] - 2024-06-25
### Changed
- 更新gotdx版本到1.22.17
- update changelog

## [1.8.31] - 2024-06-24
### Changed
- misc增加时间戳字段
- box新增时间戳字段和dxbn周期数字段
- box新增SAR指标数据
- box新增SAR数据增量计算方法
- update changelog

## [1.8.30] - 2024-06-24
### Changed
- 更新依赖库gotdx版本到1.22.16
- update changelog

## [1.8.29] - 2024-06-21
### Changed
- 更新依赖库gotdx版本到1.22.15
- update changelog

## [1.8.28] - 2024-06-20
### Changed
- 删除多余的逻辑判断
- 更新依赖库gotdx版本到1.22.14
- update changelog

## [1.8.27] - 2024-06-20
### Changed
- 更新gotdx版本到1.22.13
- update changelog

## [1.8.26] - 2024-06-19
### Changed
- 更新gotdx版本到1.22.11
- update changelog

## [1.8.25] - 2024-06-16
### Changed
- 调整版本号在开发阶段和发行后两种不同情况的处理方式
- update changelog

## [1.8.24] - 2024-06-16
### Changed
- 更新依赖库版本
- 新增实验性代码, 开发阶段获取模块版本号
- 修订主程序入口函数main的注释
- update changelog

## [1.8.23] - 2024-06-14
### Changed
- 初步实现市场情绪值的规划
- 统一冗详模式输出过滤结果
- 更新依赖库版本
- 优化板块部分代码
- 更新依赖库版本
- update changelog

## [1.8.22] - 2024-06-12
### Changed
- 删除cpmm指标,格式太旧,大智慧已经不能打开
- 调整macos的服务模板
- 调整linux服务
- 调整linux服务
- 调整linux服务
- 调整linux服务
- 暂时屏蔽linux守护进程
- 分离macos的服务模板
- 新增单一标的策略检测结果增加具体数值的尝试
- 调试Linux守护进程
- 规则新增显示冗详信息开关, 默认为false
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程执行方式
- 调整linux守护进程流程
- 优化服务service,linux系统暂时采用守护进程的方式执行
- 调整linux守护进程临时文件路径
- 规则新增显示冗详信息开关, 默认为false
- 更新依赖库
- 修订history特征组合
- update changelog

## [1.8.21] - 2024-06-06
### Changed
- 回测新增按照实际交易天数的统计数据
- 增补1号策略的注册
- 添加cpmm公式
- update changelog

## [1.8.20] - 2024-06-02
### Changed
- 修订README,修复数据命令的用法
- update changelog

## [1.8.19] - 2024-06-02
### Changed
- 优化回测宽表数据日期的对齐逻辑
- 修复除权出息日期没有对齐的bug
- 调整k线测试代码
- 优化宽表数据, 刷新历史数据用以对齐前复权数据
- 历史数据空头统计调整为只计算最低价连续新低
- 更新exchange版本到0.5.6
- 修复box中当前分钟数计算错误的bug
- update changelog

## [1.8.18] - 2024-05-27
### Changed
- update changelog
- update changelog
- update changelog

## [1.8.17] - 2024-05-27
### Changed
- 明确宽表转快照时引用特征数据的日期
- 安全分接口补充原始的页面链接
- 取消交易日18:10的全部更新
- 历史数据增加开盘量
- 调整wide转snapshot的开盘量比的计算方法
- 历史特征数据增加前日收盘
- 补充计算多头排列的注释信息
- 统一回测中的日期变量
- 回测策略单一标的支持指定日期
- 回测日期对齐
- 历史数据新增跳空缺口周期和新高次数
- 修订短线新高为量价齐升
- 新增短线量价齐跌统计字段
- update changelog

## [1.8.16] - 2024-05-23
### Changed
- 预留二次排序的todo
- 补充均线动向的字段注释
- 调整测试函数名
- 修订测试代码
- 删除废弃的函数测试代码
- 历史数据新增多头排列周期数据字段
- update changelog

## [1.8.15] - 2024-05-21
### Changed
- 调整策略概要中订单类型, 从配置文件加载
- update changelog

## [1.8.14] - 2024-05-21
### Changed
- 回测输出策略概要信息
- 新增开盘价处于上一个交易日K线实体位置的判断方法
- update changelog

## [1.8.13] - 2024-05-21
### Changed
- 优化部分代码
- 调整测试代码
- update changelog

## [1.8.12] - 2024-05-20
### Changed
- 增加高开幅度阀值
- update changelog

## [1.8.11] - 2024-05-20
### Changed
- 删除废弃的代码
- 修订交易通道参数承接比字段注释
- 策略配置参数新增低开幅度阀值
- 交易时段新增是否盘前交易方法
- 更新依赖库版本, 以支持盘前交易条件判断
- update changelog

## [1.8.10] - 2024-05-18
### Changed
- 调整存储csv文件，默认为不强制刷新, 在特殊环境比如misc的刷新，可以选择强制更新。
- 修复错误的测试代码
- 对齐misc中F10的日期为特征数据日期
- 修复订单状态文件名称不一致的bug
- 调整计算可用资金量的函数, 允许临时调整策略可交易标的总数
- 调整订单状态功能源文件名
- 拆分股票池合并功能代码
- 微调股票池合并函数, 调整参数名, 删除废弃的代码
- 优化策略单一标的可用资金量的计算方法，使分配给策略的资金尽可能足额买入。非实时订单的策略的可交易标的列表需要一次性写入股票池, 计算可用金额会用实际发生交易的股票标的数量来计算。
- 新增无效订单常量
- 新增最小订单数常量
- 修复策略订单配额判断不精准的bug, 遗漏了实时订单的最大数判断逻辑
- 股票池字段OrderStatus增加注释, 明确有效的取值及含义
- 股票池中未配置交易参数的情况输出错误日志
- 修订流程的注释, 校对交易日期比矫正日期描述准确一些, 毕竟不是发生了错误, 只是需要统一格式
- 控制台增加策略无权限提示
- 允许策略覆盖
- update changelog

## [1.8.9] - 2024-05-15
### Changed
- 修订1号策略的测试代码
- update changelog

## [1.8.8] - 2024-05-15
### Changed
- 剔除history特征中废弃的字段
- 调整快照缓存路径
- update changelog
- 剔除history特征中废弃的字段

## [1.8.7] - 2024-05-13
### Changed
- 修复快照参数重新赋值的bug
- update changelog

## [1.8.6] - 2024-05-11
### Changed
- 更新依赖库gotdx版本到1.22.6,修复服务器连接池可能被耗尽的bug
- 更新依赖库pandas版本到1.4.7
- 配置data增加快照snapshot的并发数控制, 默认为0的情况下，并发数为服务器总数的一半
- 增加快照默认最小并发数的常量
- update changelog

## [1.8.5] - 2024-05-09
### Changed
- 调整K线相对位置的判断方法, K线实体不能重合
- 删除废弃的代码
- update changelog

## [1.8.4] - 2024-05-02
### Changed
- 117跳空低开优先于固定收益率卖出逻辑
- update changelog

## [1.8.3] - 2024-05-02
### Changed
- 117增加跳空低开卖出逻辑
- update changelog

## [1.8.2] - 2024-05-02
### Changed
- 增补部分交易逻辑的注释
- 调整账户的滑动初始化组件为RollingOnce
- 修订策略配置中固定收益率FixedYield的说明, 针对买入策略执行, 会存在比较复杂的情况，所以暂时只和卖出策略的关联有效
- 修订EvaluatePriceForSell测试代码
- 117增加固定收益率的卖出逻辑
- update changelog

## [1.8.1] - 2024-04-29
### Changed
- 回测中策略的订单标识类型以配置文件为准
- backtesting配置项增加beta和alpha计算时的参考标的, 默认是上证指数
- 剔除v1和v2两个版本的回测，统一以v3版本的为主，收敛beta和alpha的计算
- update changelog

## [1.8.0] - 2024-04-28
### Changed
- 调整ism"涨停"相关字段拼音缩写的拼写错误
- 修订私有函数评估价格笼子卖出价格函数的错误注释
- 策略参数增加固定收益率, 默认是0
- 新增按照固定收益率评估卖出价格的函数
- update changelog
- 调整tail类型策略的评估价格为次日收盘价

## [1.7.9] - 2024-04-21
### Changed
- 公告风险关键词新增"立案"
- 回测增加保底的隔日溢价率
- 修复季报、年报风险提示周期的bug, 披露日期有过期已公布的情况
- 增加K线实体相对位置的判断方法, 意图类似于跳空
- 调整回测中保底隔日收益率的逻辑
- 调整快照注释和部分方法名
- 修订K线实体的高低位置
- 交易数据新增竞价承接比
- 交易数据新增竞价承接比的应用强弱判断
- 快照数据新增竞价承接比的应用强弱判断
- update changelog

## [1.7.8] - 2024-04-19
### Changed
- 回测快照增加下一个交易日的开盘、收盘、最高、最低价，便于适应各种策略的溢价计算
- 更新依赖库gotdx版本到1.22.3
- data配置项增加一年期存款利率, 默认值1.65%
- 修订波浪数据字段注释
- 调整配置data下的回测参数
- 交易参数增加无风险利率
- 调整部分代码
- 回测数据样本增加beta和alpha值
- update changelog

## [1.7.7] - 2024-04-14
### Changed
- 调整backtesting变量顺序
- !3 #I901R3 修订snapshot以满足engine的回测功能
* 修订回测中的快照数据
* 修订回测中的快照数据
- 修订回测日期存在策略没有命中情况，导致求平均值为NaN的bug
- update changelog

## [1.7.6] - 2024-04-12
### Changed
- 更新依赖库gotdx版本到1.22.2
- 基础过滤规则新增融资余额占比过大的判断
- 修订未匹配量不准确的bug, 在竞价结束后不能更新未匹配量
- update changelog

## [1.7.5] - 2024-04-11
### Changed
- 新增情绪大师(统计类)特征数据
- 修复F10数据时, 如果安全分是0, 则更新, 安全分数据没有历史数据
- 修订基础数据集缓存变量, 增加前缀__
- misc增加融资余额占比字段
- update changelog

## [1.7.4] - 2024-04-10
### Changed
- 修复配置文件结构中tag的bug,应该是yaml
- 调整财报预警周期默认值为3个交易日
- 修复没有年报披露和季报披露死循环的bug
- 临时增加更新F10数据时告警日志, 以跟踪新增数据项可能潜在的处理bug
- 过滤规则新增安全分、每股收益、每股净资产和财报披露前夕的风险期
- update changelog

## [1.7.3] - 2024-04-10
### Changed
- 修复财报披露日期检索中披露和预披露的bug
- 调整判断财报披露前风险周期的处理逻辑
- update changelog

## [1.7.2] - 2024-04-10
### Changed
- 补充波浪字段的注释
- 更新依赖库版本
- F10新增两融标注和年报季报的披露日期
- update changelog

## [1.7.1] - 2024-04-04
### Changed
- 更新依赖库版本
- 调整波浪使用字段的配置字段名成, 从filed改为fields
- 特征配置十字星在K线的占比, 默认0.5%且K线实体在上影线和下影线三者之中最小
- 修订wave字段注释
- 修订wave字段注释
- 快照新增跳空缺口判断的方法
- 删除部分废弃的代码
- 调整两融测试代码
- 更新依赖库版本
- 新增获取股价的小数点位数
- update changelog

## [1.7.0] - 2024-03-28
### Changed
- windows服务描述信息增加版本号
- 修订dfcf包名
- 新增财报数据营业总收入错误信息
- 特征配置新增波浪(wave)使用的字段及周期数
- update changelog

## [1.6.9] - 2024-03-21
### Changed
- 修复涨幅错误信息的描述
- 数据配置新增特征采用规则, 股价主导还是趋势主导
- 修订排序规则
- 如果控制台传入参数小于等于0,则使用策略配置的股票数量
- 修复黑白名单不能清空的bug
- 更新依赖库版本
- 调整快照更新时间为竞价数据结束后
- update changelog

## [1.6.8] - 2024-03-19
### Changed
- 去掉过期订单的日志
- update changelog

## [1.6.7] - 2024-03-19
### Changed
- 细化买入委托的处理过程日志
- 细化板块过滤规则, 盘中过滤不宜使用早盘数据
- update changelog

## [1.6.6] - 2024-03-19
### Changed
- 板块排序增加早盘和盘中排序两个方式, 早盘和盘中通过规则flag来判断
- 更新gotdx版本到1.21.7, 更新行业板块数据
- update changelog

## [1.6.5] - 2024-03-18
### Changed
- 更新exchange,gotdx版本,调整尾盘集合竞价的数据结束时间,预留30s给misc更新收盘数据
- update changelog

## [1.6.4] - 2024-03-17
### Changed
- 宽表增加校验和方法
- 宽表数据中, 在约定的可更新成交数据的时间范围内, 如果宽表的成交数据部分校验和为0会强制更新
- update changelog

## [1.6.3] - 2024-03-17
### Changed
- 修订注释
- 新增输出周线(非dataframe)函数
- update changelog

## [1.6.2] - 2024-03-17
### Changed
- 调整输出宽表接口
- 屏蔽曲线回归测试
- 修复motd信息错误的bug
- 更新依赖库版本
- update changelog

## [1.6.1] - 2024-03-17
### Changed
- 修订测试代码
- 调整成交量比的范围,这里的比例是当日成交量除以前日成交量,非净增比例
- 标注部分公开函数为不推荐
- 增加涨幅规则, 适用于盘中或者尾盘策略的过滤
- 剔除短期用不到的依赖库
- 修订测试代码中的bug
- 删除废弃的变量
- K线新增时间戳Datetime字段
- 拆分前复权计算方法为独立的私有函数
- 梳理更新K线的业务流程
- 修订规则中涨幅字段ChangeRate的描述
- update changelog

## [1.6.0] - 2024-03-12
### Changed
- 暂时屏蔽切换日期清理misc内存的操作
- 复原清理隔日内存
- 更新依赖库版本及go版本
- update changelog

## [1.5.9] - 2024-03-12
### Changed
- 更新依赖库版本
- update changelog

## [1.5.8] - 2024-03-11
### Changed
- 更新依赖库版本
- update changelog

## [1.5.7] - 2024-03-11
### Changed
- 调整默认的最大放量比例为1.8倍
- update changelog

## [1.5.6] - 2024-03-11
### Changed
- 增加misc数据获取失败的错误日志
- update changelog

## [1.5.5] - 2024-03-10
### Changed
- 更新依赖库num,pandas,ta-lib版本
- 更新依赖库版本
- 优化内存显示
- 调整捕获异常的参数
- 使用//go:linkname从internal/cpu包中导出获取cpu型号函数
- 调整测试代码
- 更新依赖库版本
- update changelog

## [1.5.4] - 2024-02-25
### Changed
- 更新依赖库版本
- 更新依赖库num版本到0.1.2
- 更新依赖库版本
- update changelog

## [1.5.3] - 2024-02-19
### Changed
- 调整pandas接口
- 更新依赖库版本
- update changelog

## [1.5.2] - 2024-02-12
### Changed
- 更新ta-lib版本
- update changelog

## [1.5.1] - 2024-02-12
### Changed
- 修订初始化接口Init方法注释, 去掉证券代码
- 修订miniQMT配置中订单路径的注释
- 修订资金范围的错误信息
- 增加开发可能遇到问题的解决方案。 可能会出现"call has possible Printf formatting directive %s"的问题,这是由于go vet导致的, 函数本身没有问题。
- 迁移成交数据常量到exchange
- 适配新版本的pandas
- update changelog

## [1.5.0] - 2024-01-31
### Changed
- 调整修复操作时进度条不满100%的bug
- 实验misc新因子
- 新增多空趋势信号及周期
- 修复宽表更新时由于成交数据只更新当日数据造成的缺少上一个交易数据的bug
- update changelog

## [1.4.9] - 2024-01-30
### Changed
- 修复repair命令行参数日期范围可能存在休市的情况
- update changelog

## [1.4.8] - 2024-01-30
### Changed
- 调整misc结构体名, 去掉早期的exchange字样
- 修复宽表计算的bug
- update changelog

## [1.4.7] - 2024-01-28
### Changed
- 更新依赖库pandas,ta-lib版本
- 调整ma源文件名
- 新增实时计算EMA的函数
- 新增实时计算MACD的函数
- update changelog

## [1.4.6] - 2024-01-27
### Changed
- 默认配置文件增加订单路径
- 更新依赖库版本
- update changelog

## [1.4.5] - 2024-01-25
### Changed
- 调整进度条, 增加wait方法
- update changelog

## [1.4.4] - 2024-01-25
### Changed
- 更新依赖库版本
- update changelog

## [1.4.3] - 2024-01-25
### Changed
- 更新依赖库版本
- 新增每日8:55更新网络定时任务
- update changelog

## [1.4.2] - 2024-01-24
### Changed
- 更新依赖库版本
- update changelog

## [1.4.1] - 2024-01-24
### Changed
- 修复协程方式拉取快照引起的misc和snapshot map数据读写竞争的bug
- update changelog

## [1.4.0] - 2024-01-23
### Changed
- 更新依赖库gotdx版本
- update changelog

## [1.3.9] - 2024-01-23
### Changed
- 更新依赖库gotdx版本
- update changelog

## [1.3.8] - 2024-01-23
### Changed
- 更新gox库, 屏蔽cache.FastCache代码
- update changelog

## [1.3.7] - 2024-01-23
### Changed
- 拆分出进度条bar的index
- 调整进度条bar的index
- 调整进度条归属源文件到command_update
- 删除废弃的代码
- 更新依赖库gotdx版本
- 更新依赖库gotdx版本, 强化panic之前记录日志
- update changelog

## [1.3.6] - 2024-01-22
### Changed
- 更新依赖库exchange,gotdx版本号
- update changelog

## [1.3.5] - 2024-01-22
### Changed
- 调整编译脚本允许项目目录外执行脚本
- 新增mac操作系统arm64编译脚本
- 新增windows操作系统arm64编译脚本
- 调整mac arm64编译脚本
- 调整windows arm64编译脚本
- 新增linux amd64编译脚本
- 新增windows下ps1编译脚本
- 新增windows下python的编译脚本
- update changelog

## [1.3.4] - 2024-01-22
### Changed
- 修订实时更新快照的队列长度不超过服务器数量的一半
- update changelog

## [1.3.3] - 2024-01-22
### Changed
- 更新依赖库版本
- update changelog

## [1.3.2] - 2024-01-22
### Changed
- 优化部分缓存代码
- 拆分每日订单缓存处理方法
- 调整股票池结构
- 简化更新股票池功能, 取消单一策略的选股列表的落盘, 统一归到股票池
- 优化持仓部分代码
- 重做持仓周期列表
- 优化每日全部更新的完成状态文件的处理方法
- 优化本地订单文件的日期列表函数
- 优化订单日期的列表, 从本地缓存的订单文件列表获取
- 更新依赖库版本
- 修订tracker有效策略为空时, 应该输出info类型日志
- 调整记分牌String方法
- 更新gotdx版本号
- 执行策略之前判断是否配置了交易规则
- 增加策略执行时个股列表的三个来源的描述
- 简化实时更新任务的日志
- 适配exchange
- 收敛日期字符串到uint32类型转换的函数到exchange
- !2 snapshot多线程
* 多线程更新snapshot
- 更新依赖库版本
- 补充公开函数的注释
- 控制台输出统计表格, 去掉涨速和力度两个字段
- update changelog

## [1.3.1] - 2024-01-17
### Changed
- 调整宽表部分信息
- 宽表的获取缓存文件, 暂时不做内存缓存
- 调整测试代码
- 修复持仓存在上一个交易日卖出未成交的持仓bug, 修正的逻辑应该继续当作最后一个持股日卖出
- update changelog

## [1.3.0] - 2024-01-17
### Changed
- 修复宽表重做丢失昨日收盘和涨跌幅的bug
- 更新时如果自由流通股本为0, 则用流通股本覆盖
- 调整记分牌结构, 增加名称
- 更新gox版本号, 进度条增加等待结束信号
- 调整数据适配器部分结构名
- update changelog

## [1.2.9] - 2024-01-16
### Changed
- 修订部分即将废弃的函数
- 更新gotdx版本
- 清楚过时的todo
- 新增宽表的缓存机制
- update changelog

## [1.2.8] - 2024-01-15
### Changed
- 调整命令字, 以支持服务名等关键信息, 可以从下游的项目名传递过来
- update changelog

## [1.2.7] - 2024-01-15
### Changed
- 调整成交记录的取法, 先从缓存获取, 缓存没有再从服务器读取
- update changelog

## [1.2.6] - 2024-01-15
### Changed
- 删除废弃的日志输出
- 删除废弃的日志输出
- 移除并修复收盘判断当日的bug
- K线数据少于约定的120天会提示没有K线
- update changelog

## [1.2.5] - 2024-01-15
### Changed
- 权限检测失败前, 不显示欢迎语
- 删除不必要的日志
- 调整misc实时更新的任务名称
- 调整tracker部分输出信息
- mac笔记本misc实时更新存在文件名空的现象, 但是其它机器未发现, 增加一条文件名切换时的告警信息输出
- update changelog

## [1.2.4] - 2024-01-15
### Changed
- 调整欢迎语的显示顺序
- update changelog

## [1.2.3] - 2024-01-15
### Changed
- 调整策略编号的数据类型为uint64, 修改权限注册接口入参为uint64, 和权限模块保持一致
- update changelog

## [1.2.2] - 2024-01-15
### Changed
- 更新依赖库版本
- 调整日期范围函数
- 从zero-sum迁移策略执行函数
- 命令行统一增加欢迎语
- 移除废弃的代码
- 优化策略执行方式
- 增加盘中扫描前策略有效性确认
- update changelog

## [1.2.1] - 2024-01-14
### Changed
- 修订注释描述错误的问题
- 更新gox版本, 以支持recover时, 允许传入可变参数, 提供更多的造成异常的证据
- 更新gox版本, 去掉无异常时的无意义解析入参
- 更新gox版本, 优化可以忽略panic的用法
- 更新exchange版本, 强化交易日期范围的检查
- 收敛panic的处理方式, 统一归于gox的runtime
- update changelog

## [1.2.0] - 2024-01-13
### Changed
- 更新依赖库版本
- tracker增加校验策略权限
- engine增加tracker模块及策略权限验证
- 修复wide可能存在的数据日期错乱的情况
- update changelog

## [1.1.9] - 2024-01-13
### Changed
- 修复了日期范围函数因为前后日期颠倒引发的异常, 新增校验wide和k线开始日期是否对齐
- 更新依赖库版本
- 更新依赖库版本
- update changelog

## [1.1.8] - 2024-01-13
### Changed
- 修订部分代码的变量名, 结构体名
- 修复命令参数错误时控制台输出两次错误信息的bug
- 调整历史成交数据的部分函数名, 收敛关于历史成交记录的默认日期
- update changelog

## [1.1.7] - 2024-01-12
### Changed
- 整理部分代码, 合并小功能
- update changelog

## [1.1.6] - 2024-01-12
### Changed
- 新增板块测试代码
- 从stock迁移扫描功能
- update changelog

## [1.1.5] - 2024-01-12
### Changed
- 更新exchange版本号
- 梳理数据集代码
- 保存旧版本wide的更新操作方法
- 调整宽表、历史成交数据的更新方案
- update changelog

## [1.1.4] - 2024-01-11
### Changed
- 调整misc缓存文件名
- 更新pkg版本, 修复csv加载不能解析从科学计数法的浮点转换成int64的bug
- 调整宽表的结构
- update changelog

## [1.1.3] - 2024-01-11
### Changed
- 适配exchange工具包
- update changelog

## [1.1.2] - 2024-01-11
### Changed
- exchange增加买入和卖出金额
- update changelog

## [1.1.1] - 2024-01-11
### Changed
- 修订换手z的算法, 板块指数类的自由流通股本用流通股本来计算
- 修订换手z的算法, 板块指数类的自由流通股本用流通股本来计算
- 修复买入方向的算法
- update changelog

## [1.1.0] - 2024-01-10
### Changed
- 拆分每日系统初始化功能函数
- 调整9点整重置定时任务的key为global_reset, clean的字面意义已过时
- 调整流通股本默认最小值0.5亿股
- 更新依赖库版本
- update changelog

## [1.0.9] - 2024-01-03
### Changed
- 修复f10季报死锁的bug
- update changelog

## [1.0.8] - 2024-01-01
### Changed
- 更新gotdx, 修复扩展行情和标准行情相同的bug
- update changelog

## [1.0.7] - 2024-01-01
### Changed
- 调整errors的用法
- 更新依赖库gotdx版本
- 迁移fastjson到pkg
- 修复特征history名字错误的bug
- 优化update和repair子命令
- 测试goland在git提交时格式化的问题
- 恢复测试前的代码
- 更新gotdx依赖库版本
- 调整买入价格逻辑, 从当前价格默认增加0.05
- 测试代码: 新增,批量转换证券代码为qmt支持的格式
- 消除部分未使用变量的告警信息
- 优化cache1d的内存申请方式
- crontab配置去掉name字段
- 策略新增关于价格笼子的两个参数
- 优化调度任务, 修订同种类型定时任务的源文件名前缀
- 优化账户数据
- 为每一组qmt常量设定类型
- 优化持仓信息
- 委托买入启用价格笼子算法, 价格笼子相关的参数从配置文件加载
- 调整qmt订单路径, 修复不能从配置文件读取的bug
- 新增定时任务, 每个交易日15点02分同步全天的委托订单
- update changelog

## [1.0.6] - 2023-12-30
### Changed
- 修复按策略数量分摊可用资金的算法, 最后一个策略获得剩余的全部资金
- 修订部分代码
- 更新gotdx版本,修复交易日历用到最后1天时无法引用下一个交易日的bug
- update changelog

## [1.0.5] - 2023-12-30
### Changed
- 更新依赖库版本
- 修复策略接口源文件拼写的错误
- 优化模型部分配置性代码
- 调整url常量前缀, 去掉k
- 删除废弃的代码
- 测试通达信问小达的接口, 失败
- 预留增加两融详细数据
- 更新gotdx版本
- 增加东方财富两融数据接口的实现
- update changelog

## [1.0.4] - 2023-12-29
### Changed
- 新增工具集子命令, 实现了tail功能
- update changelog

## [1.0.3] - 2023-12-28
### Changed
- 调整跳空低开的开关到规则参数结构中
- 跳空低开统一加在基础规则中
- update changelog

## [1.0.2] - 2023-12-28
### Changed
- 修订测试代码
- 调整测试代码
- 增加检查指定个股在策略中的执行情况
- 回测增加check个股在策略中的过滤情况, 输出失败的详细信息
- update changelog

## [1.0.1] - 2023-12-27
### Changed
- 修订默认配置
- 修订默认配置
- 调整策略参数
- update changelog

## [1.0.0] - 2023-12-26
### Changed
- NumberRange调整最大最小值的字段名, 并新增获取数值的方法
- 优化卖出规则
- 配置项第一层去掉rules, 该在每一个策略中配置
- 调整过滤规则, 规则参数从策略中加载
- update changelog

## [0.9.9] - 2023-12-25
### Changed
- 调整可用金额的计算方法, 以当日可用为主, 不为下一交易日做预留处理
- update changelog

## [0.9.8] - 2023-12-25
### Changed
- 更新gotdx库,增加更多的服务器IP地址
- 命令字增加检测服务器网速
- update changelog

## [0.9.7] - 2023-12-24
### Changed
- 更新ta-lib版本
- 调整数值范围的验证逻辑, 如果begin和end都为0, 视为不验证, 默认通过
- 确定出了流通股本和流通市值以外, 其它默认都是不验证
- 调整topn的取法
- 策略接口的filter方法增加参数交易规则
- 配置文件删除order项
- 拆分个股列表, 增加过滤功能
- update changelog

## [0.9.6] - 2023-12-24
### Changed
- 适配gox新版本的http函数
- 删除protobuf的测试代码, 放弃rpc的想法
- 删除部分废弃的代码
- 调整特征源文件名, 保持前缀feature
- update changelog

## [0.9.5] - 2023-12-22
### Changed
- 删除废弃的科创板过滤规则
- 调整过滤规则注释中的逻辑序号
- update changelog

## [0.9.4] - 2023-12-22
### Changed
- 修复每次解析文本是session数组没有重置的bug
- update changelog

## [0.9.3] - 2023-12-21
### Changed
- 测试增加一刀切获取持股到期日的个股列表
- 交易规则增加买入是否支持跳空低开
- 默认允许跳空低开
- 交易规则中板块列表支持前缀带减号, 表明要剔除的板块成分股
- update changelog

## [0.9.2] - 2023-12-21
### Changed
- 修复数字范围间隔符号不能用-的bug
- update changelog

## [0.9.1] - 2023-12-21
### Changed
- 修订订单备注确实的问题
- 调整规则配置, 增加NumberRange功能
- 调整规则配置, 增加NumberRange功能
- update changelog

## [0.9.0] - 2023-12-20
### Changed
- 交易规则增加可止盈和可止损的判断方法
- update changelog

## [0.8.9] - 2023-12-20
### Changed
- 调整ta-lib版本号
- update changelog

## [0.8.8] - 2023-12-20
### Changed
- 调整数据集源文件名
- 统一数据到factors目录
- update changelog

## [0.8.7] - 2023-12-20
### Changed
- 清除废弃的代码
- 调整box部分字段注释
- 调整规则中市值的范围
- 交易规则增加止盈止损
- 调整box的字段注释
- 调整交易规则中的止盈字段字段名
- 调整117号卖出策略的逻辑
- update changelog

## [0.8.6] - 2023-12-19
### Changed
- 调整engine内置的特征组合文件名
- 开放exchange和box特征数据
- update changelog

## [0.8.5] - 2023-12-19
### Changed
- 整理部分代码
- 统一engine中关于行情快照结构体的定义
- update changelog

## [0.8.4] - 2023-12-18
### Changed
- 修订测试策略的接口实现
- 调整最小市值默认值为5亿
- F10增加市值控制
- update changelog

## [0.8.3] - 2023-12-18
### Changed
- history增加昨日开盘,收盘,成交量和成交额.
- update changelog

## [0.8.2] - 2023-12-18
### Changed
- 去除多余的echo
- 修订编译脚本, 从go.mod中读取module
- 更新依赖库版本, gotdx的snapshot增加本地时间戳字段, 用以观察本地时钟和服务器的差距
- 简化一刀切卖出规则
- 删除废弃的代码
- 买入和卖出检查是否黑白名单
- 更新依赖库版本
- 删除废弃的代码
- 优化部分代码
- update changelog

## [0.8.1] - 2023-12-17
### Changed
- 增加3个可能用到的时间戳格式
- 移除时间戳格式
- 特征数据增加验证样本的方法
- 性能分析默认调整为关闭状态
- 增加交易时段测试代码
- 调整series函数
- 新增命令行永久flag, avx2加速和cpu核数控制
- 增加使用go build -ldflags构建时传入的版本号的提示性注释
- 调整策略Strategy的评估接口, result参数改用泛型treemap
- history增加9和19日的均价、均量线
- 调整获取最新数据的行数,明确数据类型
- update changelog

## [0.8.0] - 2023-12-16
### Changed
- 增加文档
- 新增黑白名单功能
- update changelog

## [0.7.9] - 2023-12-15
### Changed
- 更新pandas版本号
- 修复code list可能存在变化的情况引发进度条异常的bug
- 新增通过策略编码获取用于qmt系统的string类型的策略名, 大写S作为前缀后面跟quant1x系统的策略编码
- 交易规则增加持股周期字段
- 交易时段新增判断是否当天的最后一个交易时段
- 新增qmt相关的功能函数
- 调整去重函数
- history增加2日均线和4日均线
- 优化和调整交易规则结构
- 优化qmt策略名称的处理方法
- 新增拉取指定日期内策略首次买入的个股列表
- 新增价格笼子的计算方法
- 更新依赖库版本
- 删除废弃的代码
- 合并部分小功能代码
- 增加增量计算均线的函数
- 细化卖出策略
- 调整卖出定时任务
- 收敛engine中的所有command, 目的是为了统一处理debug开关
- 拆分出runtime配置项, runtime参数增加debug开关, 默认关闭
- 修订最新的配置项的默认配置文件
- 修订配置文件的处理方法
- 修订一刀切的常量名, 用Sell替代Sale
- 优化配置加载过程
- 修订卖出策略中订单备注的信息
- 修订卖出策略中订单备注的信息
- 更新依赖库版本
- 删除独立的卖出策略配置sell
- 拟增加各Flag订单的总开关
- 新增861号卖出策略编码ID的常量
- 新增漏掉的卖出时段判断
- 屏蔽废弃的功能函数
- 服务方式运行更新快照去掉进度条
- 优化获取当前交易日期的方法
- update changelog

## [0.7.8] - 2023-12-13
### Changed
- 使用go 1.21.5
- 增加doc文档说明性源文件
- 更新依赖库版本
- 收敛获取应用程序文件名的方法
- 执行策略前输出策略概要
- 调整util工具包
- update changelog

## [0.7.7] - 2023-12-12
### Changed
- 更新依赖库pkg, ta-lib版本号
- 调整print指令, 增加输出缓存日期和特征日期
- F10增加财务数据报告期
- 订单字段增加tag式注释
- 增加qmt的常量
- 增加持仓结构
- 调整qmt持仓字段
- 交易员参数增加交易角色
- 调整交易参数
- 更新依赖库gox版本号
- 调整订单结构
- 规则增加过滤股票代码前缀
- 交易配置增加费率
- 新增交易费用结构体
- 拆分出交易配置对象
- 增加通过预算输出交易费用对象
- 优化买入卖出交易费用的计算方法
- 增加mac和windows平台的amd64编译搅脚本
- 调整分割线
- 优化编译脚本
- 增加账户和策略可用资金的计算方法
- 统一最低要求K线数量的常量为120
- 特征数据增加K线数据的最低要求限制的检查
- 增加进入股票后直接向qmt proxy发起委托下单
- 更新依赖库版本
- 新增通过配置调整定时任务的开关和触发条件
- 调整策略的错误信息
- 修订规则的错误信息
- 调整定时任务的配置方式
- 删除废弃的时间类测试代码
- 修订数据适配器相关的错误信息
- 交易员参数增加账户ID
- 拆分出股票代码列表的函数
- 新增配置测试代码
- 新增矫正策略交易时段的处理
- 基础数据, 拟增加两融标的
- 修订TODO注释
- 按策略关联板块以及是否过滤两融
- 优化策略有效性判断
- 更新gotdx版本
- 调整订单状态机
- 策略方面, 增加一个特殊的卖出策略117(一刀切), 新增QMT体系中的策略名函数和订单备注函数
- 调整测试代码
- 拟增加持仓列表功能
- 新增计算涨跌停板价格的函数
- 更新依赖库版本
- 拆分下单委托函数, 支持直接填充策略名和订单备注
- 实现一刀切的功能
- update changelog

## [0.7.6] - 2023-12-05
### Changed
- 更新依赖库pandas版本号
- 调整控制的最大订单数
- 删除废弃的评估方法
- 修复错误的注释
- 调整安全分http请求的方法
- 调整http get请求的参数
- 更新gox版本, http增加post方法支持
- 删除废弃的交易参数结构体
- 增加查询委托功能
- 增加委托和撤单两个函数
- 修订日志中的错误描述
- 修订日志中的错误描述
- F10增加营业总收入
- proxy服务器地址
- F10增加每股收益扣除字段
- 修复repair --all 特征数据不生效的bug
- 调整交易方向类型
- 更新依赖库gotdx,pkg版本号
- 增加撤单时段
- 优化配置加载方式
- 新增一个只获取一只股票tick数据的函数, 目的是为了方便单元测试
- update changelog

## [0.7.5] - 2023-12-03
### Changed
- 拟增加权限模块
- 新增交易模块
- 从cache目录中拆分出config
- 增加流通市值
- 调整目录结构
- 增加交易参数配置
- 增加市场雷达功能
- 加载配置文件增加错误日志
- 调整统计参数归于模型
- update changelog

## [0.7.4] - 2023-11-27
### Changed
- 修复订单状态被覆盖的bug
- 更新依赖库版本
- update changelog

## [0.7.3] - 2023-11-26
### Changed
- 配置文件新增定时任务开关
- 增加1号策略通达信公式源代码
- 股票池结构体增加活跃度和涨速字段
- 统计结构体增加涨速字段
- 调整部分快照相关的函数名
- 调整package
- 剥离部分runtime功能到gox
- 删除废弃的测试代码
- 废弃部分字段
- 更新gox版本
- update changelog

## [0.7.2] - 2023-11-19
### Changed
- 修订股票池
- 更新依赖库版本
- 更新依赖库版本
- 增加公开函数的注释
- 拟增加数据源切换功能
- 调整history csv字段名
- 调整no1的方法顺序
- 调整策略的执行方法, 改用注册的方式。暂时屏蔽回测功能。
- 拆分策略结果结构体为一个独立的源文件
- 收敛长期不更新的依赖库到pkg
- 实现strategy接口的1号策略, 去掉指针接收器的用法
- 删除独立的次新股规则
- 优化规则分组
- update changelog

## [0.7.1] - 2023-11-13
### Changed
- 更新gotdx版本号, 更新内置的板块数据文件
- update changelog

## [0.7.0] - 2023-11-13
### Changed
- 修复振幅最小值的key重复的bug
- update changelog

## [0.6.9] - 2023-11-13
### Changed
- 规则增加振幅范围0.00%~15.00%
- update changelog

## [0.6.8] - 2023-11-13
### Changed
- 更新依赖库版本
- 增加一个预备的投票模块
- 调整部分函数为公开
- update changelog

## [0.6.7] - 2023-11-07
### Changed
- 去除废弃的代码
- 增加定时任务配置结构体
- 增加更新内存中的K线操作
- 调整记分牌的tag
- 增加无效周期的常量
- 增加盘中实时订单标识
- 调整基础过滤规则
- 更新ta-lib版本号
- 更新其它非quant1x组织的依赖库版本
- 修复股本变化的类型中遗漏送配股上市的bug
- 更新gotdx版本, 优化除权除息的股本变化处理方法
- update changelog

## [0.6.6] - 2023-10-31
### Changed
- 更新ta-lib版本
- 调整规则和订单配置加载方式
- 升级依赖库版本号
- 去除废弃的代码
- update changelog

## [0.6.5] - 2023-10-30
### Changed
- 调整MV5的计算方法
- update changelog

## [0.6.4] - 2023-10-30
### Changed
- 调整精度条bar的序号处理方式
- update changelog

## [0.6.3] - 2023-10-30
### Changed
- 调整bar的空白行
- 数据2个空白行, 暂时先这么固定输出, 后面再优化调度任务结构
- update changelog

## [0.6.2] - 2023-10-30
### Changed
- 增加业绩预告数据
- update changelog

## [0.6.1] - 2023-10-29
### Changed
- 增加leveldb测试代码
- 调整models
- 更新依赖库版本
- 更新gotdx版本
- 分时数据增加异常捕获
- 修复map并发读写的bug
- 数据集更新增加捕获异常
- 增加策略文件缓存路径的常量
- 策略接口增加订单类型和过滤器
- 次新股默认规则通过
- 调整F10的证券名称字段
- 调度任务启动时增加互斥锁
- 屏蔽实时更新K线的定时任务
- 更新gox版本
- 删除废弃的pprof代码
- 特征数据增加异常捕获
- 增加快照定时任务
- 修复snapshot可能是nil的bug
- 更新gotdx版本
- 更新gox版本
- 新增具有滑动窗口速度控制的WaitGroup
- 增加从单个snapshot更新K线的函数
- 调整季报的处理方法
- 更新gotdx版本
- 更新依赖库版本
- 优化1d缓存对象
- history增加前5日分钟均量的方法
- F10的公告信息只处理证券代码
- 调整证券名称、季报的处理方法
- 调整GoMaxProcs函数路径
- 调整GoMaxProcs函数名
- 增加Recover函数
- 调整捕获panic的函数名
- 优化debug开关
- 收敛recover捕获panic异常的方法
- 删除废弃的重置证券代码的处理方法
- 调整股票池StockPool的package
- 调整存储订单的同时输出到股票池
- 优化imports
- update changelog

## [0.6.0] - 2023-10-25
### Changed
- 增加调度任务日志
- update changelog
- 增加调度任务日志

## [0.5.9] - 2023-10-25
### Changed
- 调整F10的csv字段
- update changelog

## [0.5.8] - 2023-10-25
### Changed
- 调整调度任务info级别日志内容
- update changelog

## [0.5.7] - 2023-10-25
### Changed
- 优化实时更新K线的时间范围
- update changelog

## [0.5.6] - 2023-10-25
### Changed
- 增加回测功能
- 修订README中各模块的完成情况
- 主程序增加回测模块
- 增加过滤规则功能
- ResourcesPath改为常量
- 增加输出规则列表的子命令
- 子命令增加测试参数异常的测试性代码
- 优化规则引擎
- update changelog

## [0.5.5] - 2023-10-24
### Changed
- 增加修订application的初始化代码
- update changelog

## [0.5.4] - 2023-10-24
### Changed
- 新增K线宽表数据的基础函数
- 完善宽表数据
- update changelog

## [0.5.3] - 2023-10-24
### Changed
- 增加分时数据缓存
- update changelog

## [0.5.2] - 2023-10-23
### Changed
- 修复xdxr缺少date和code的bug
- update changelog

## [0.5.1] - 2023-10-23
### Changed
- 调整服务接口
- 调整定时任务的回调函数
- 调整调度任务代码结构
- update changelog

## [0.5.0] - 2023-10-23
### Changed
- 调整实时更新K线的兜底逻辑
- update changelog

## [0.4.9] - 2023-10-22
### Changed
- 增加实时更新K线
- 实时更新K线增加内外盘两个字段
- 恢复非交易时段的不操作的逻辑
- 更新依赖库版本
- 调整定时任务的实现方式
- update changelog

## [0.4.8] - 2023-10-21
### Changed
- 更新gox版本
- 调整等待应用结束的机制
- update changelog

## [0.4.7] - 2023-10-21
### Changed
- 修正测试代码
- 修订业绩预报结构体注释
- engine增加性能分析工具
- 给pprof增加开关
- update changelog

## [0.4.6] - 2023-10-20
### Changed
- 基础数据增加实时更新基础K线的函数
- update changelog

## [0.4.5] - 2023-10-20
### Changed
- 清理废弃的代码
- 更新gox版本号
- 增加系统服务子命令
- update changelog

## [0.4.4] - 2023-10-19
### Changed
- 增加daemon服务命令字
- update changelog

## [0.4.3] - 2023-10-19
### Changed
- 更新依赖版本

## [0.4.2] - 2023-10-19
### Changed
- 优化代码结构
- update changelog

## [0.4.1] - 2023-10-19
### Changed
- 调整字段名
- update changelog

## [0.4.0] - 2023-10-19
### Changed
- 调整数据接口归类划分
- 调整数据接口归类划分
- update changelog

## [0.3.9] - 2023-10-19
### Changed
- 调整数据接口归类划分
- update changelog

## [0.3.8] - 2023-10-19
### Changed
- 细分数据接口
- update changelog

## [0.3.7] - 2023-10-18
### Changed
- 增加一个轻量的特性接口, 用来扩展子特征
- 调整csv字段名
- update changelog

## [0.3.6] - 2023-10-18
### Changed
- 调整manifest结构体字段
- update changelog

## [0.3.5] - 2023-10-18
### Changed
- 调整manifest结构体私有为公开
- update changelog

## [0.3.4] - 2023-10-18
### Changed
- 新增数据集和特征的manifest
- 提取抽象结构
- update changelog

## [0.3.3] - 2023-10-17
### Changed
- 修订缓存适配器接口的注释
- 调整数据适配器接口的方法顺序
- 删除废弃的代码
- 调整数据适配器接口的方法顺序
- 调整cache1d的new函数
- 收敛cache1d的缓存文件路径函数
- 修订cache1d结构体的注释
- 优化代码
- 修订上一个季报没公布导致前十大流通股东列表为空的bug, 如果未公布, 应该沿用再上一个季度的数据
- 调整数据接口
- update changelog

## [0.3.2] - 2023-10-17
### Changed
- 增加股票池结构, 所有的数据都放在一个文件里面
- 股票池增加规则字段
- 新增数据接口
- 新增summary和trait两个接口
- 新增规则接口
- 股票池增加策略状态字段
- 调整dataset方法
- 调整dataset方法
- trait特性接口增加提供者方法
- 新增 数据的控制台命令支持接口
- 调整提供者的方法名
- 调整记分牌的package
- 增加忽略pprof文件
- 增加数据运算接口
- 增加数据项接口
- 更新gox版本
- 应用程序增加性能分析功能
- 调整数据接口
- 调整增量(不推荐)接口的package
- 调整cache1d的缓存路径
- 调整F10的csv字段名
- 调整HousNo1的csv字段名
- 调整history结构的csv字段名
- 修订缓存操作接口的注释
- 修订项目的主要关键词解释
- 更新主要依赖库版本
- aaa
- 调整代码结构
- 调整源代码文件名
- update changelog

## [0.3.1] - 2023-10-13
### Changed
- 增加ants协程池控制并发数量
- update changelog

## [0.3.0] - 2023-10-13
### Changed
- 测试协程方式跑特征数据
- update changelog

## [0.2.9] - 2023-10-13
### Changed
- 优化update和repair数据处理流程
- update changelog

## [0.2.8] - 2023-10-13
### Changed
- 优化update和repair数据处理流程
- update changelog

## [0.2.7] - 2023-10-13
### Changed
- 优化update和repair数据处理流程
- update changelog

## [0.2.6] - 2023-10-12
### Changed
- 增加周线,月线函数
- update changelog

## [0.2.5] - 2023-10-12
### Changed
- 调整engine数据的提供者为engine
- update changelog

## [0.2.4] - 2023-10-12
### Changed
- 命令字初始化改为显式
- update changelog

## [0.2.3] - 2023-10-12
### Changed
- 修改错误名
- 调整历史成交记录的update和repair, 更新的日期应该采用cacheDate
- 变更源文件名
- 增加通达信F10的资金流向, 这个数据因为网络请求的轮询机制, 数据很有可能存在不同源的问题, 从而导致数据不完整或者不正确
- 增加一个单独的增量计算的接口备用
- 新增通达信自选股列表导出函数
- 调整缓存的工厂用法
- 新增数据验证check接口
- 调整子命令的检索逻辑
- update changelog

## [0.2.2] - 2023-10-11
### Changed
- 更新gotdx版本, 历史成交数据去掉用pandas的方式读写, 改为切片和csv文件直接交换
- update changelog

## [0.2.1] - 2023-10-11
### Changed
- 修订切片自动扩容地址变化引起的优先级较高的特征信息不能打印的bug
- update changelog

## [0.2.0] - 2023-10-11
### Changed
- 增加位图, 为将来扩展特征类型做准备
- 数据插件增加get接口
- 基础数据增加历史成交数据
- 调整进度条的index
- 调整源文件名
- 屏蔽暂时废弃的变量声明
- 修订bitmap结构体注释
- 子命令print自动检测是否打印特征数据, 暂时不支持结构嵌套
- update changelog

## [0.1.9] - 2023-10-10
### Changed
- 将内部函数公开
- update changelog

## [0.1.8] - 2023-10-10
### Changed
- 增加注释
- 标注废弃部分函数
- 消除没有使用参数的告警提示
- 更新gox版本
- 调整更新和修复子命令
- update changelog

## [0.1.7] - 2023-10-10
### Changed
- 更新F10中公告的增持和减持的字段名
- 删除废弃的测试代码
- 修正安全分单词
- 增加smart接口
- 增加插件接口, 用以收盘写数据操作
- 调整插件接口名
- 移除测试性代码
- 收录github.com/mattn/go-runewidth@v0.0.15
- 修订README, 增加对于协同开发方面的说明
- 调整插件模式的遍历方法
- 调整变量的写法
- 调整变量的写法
- 调整变量的写法
- 调整变量的写法
- 调整变量的写法
- 修订分支的描述
- 修订项目总名称
- 调整基本面数据的优先级
- update changelog

## [0.1.6] - 2023-10-08
### Changed
- 删除废弃的特征组合box
- update changelog

## [0.1.5] - 2023-10-08
### Changed
- 更新gox库版本
- 调整代码归属
- 增补规范的文件名函数
- 修正cache1d的缓存关键字
- 增加个股安全评估数据
- 新增F10基本面特征数据组合
- 增加异常是显示调用栈
- repair增加基础数据
- repair增加特征数据
- 特征增加侯总1号策略
- 增加通达信协议日期转换函数
- 调整除权除息列表的测试代码
- 新增东方财富数据的接口
- 调整数据集合, 增加基础K线, 财报, 安全分, 除权除息
- 更新依赖库的版本
- 优化命令行参数解析
- 增加version, print子命令
- 调整测试代码
- 调整缓存机制的时间函数的package归属
- update changelog

## [0.1.4] - 2023-10-07
### Changed
- 更新gox、gotdx库版本
- update changelog

## [0.1.3] - 2023-10-06
### Changed
- 调整基础数据集合
- 拆分dataset
- 更新gox版本
- 调整策略结果结构体
- 调整策略结果结构体字段顺序
- 执行策略之前增加同步即时行情数据的过程, 以便策略可以使用增量计算方法
- 调整数据集和特征组合
- update changelog

## [0.1.2] - 2023-10-02
### Changed
- 完成第一个策略演示
- update changelog

## [0.1.1] - 2023-10-01
### Changed
- 新增K线和除权除息的基础数据
- 增加趋势反转代码
- add ChangeLog
- 增加第一个策略执行的demo
- update changelog

## [0.1.0] - 2023-09-28
### Changed
- first commit
- add LICENSE.

Signed-off-by: 王布衣 <wangfengxy@sina.cn>
- 增加统一的常量模块
- 新增基础k线测试程序
- 新增历史数据结构
- 新增快照数据结构
- history增加日期的描述
- 修订README
- 新增测试特征接口的代码, 以日K线为样本


[Unreleased]: https://gitee.com/quant1x/engine.git/compare/v1.9.0...HEAD
[1.9.0]: https://gitee.com/quant1x/engine.git/compare/v1.8.46...v1.9.0
[1.8.46]: https://gitee.com/quant1x/engine.git/compare/v1.8.45...v1.8.46
[1.8.45]: https://gitee.com/quant1x/engine.git/compare/v1.8.44...v1.8.45
[1.8.44]: https://gitee.com/quant1x/engine.git/compare/v1.8.43...v1.8.44
[1.8.43]: https://gitee.com/quant1x/engine.git/compare/v1.8.42...v1.8.43
[1.8.42]: https://gitee.com/quant1x/engine.git/compare/v1.8.41...v1.8.42
[1.8.41]: https://gitee.com/quant1x/engine.git/compare/v1.8.40...v1.8.41
[1.8.40]: https://gitee.com/quant1x/engine.git/compare/v1.8.39...v1.8.40
[1.8.39]: https://gitee.com/quant1x/engine.git/compare/v1.8.38...v1.8.39
[1.8.38]: https://gitee.com/quant1x/engine.git/compare/v1.8.37...v1.8.38
[1.8.37]: https://gitee.com/quant1x/engine.git/compare/v1.8.36...v1.8.37
[1.8.36]: https://gitee.com/quant1x/engine.git/compare/v1.8.35...v1.8.36
[1.8.35]: https://gitee.com/quant1x/engine.git/compare/v1.8.34...v1.8.35
[1.8.34]: https://gitee.com/quant1x/engine.git/compare/v1.8.33...v1.8.34
[1.8.33]: https://gitee.com/quant1x/engine.git/compare/v1.8.32...v1.8.33
[1.8.32]: https://gitee.com/quant1x/engine.git/compare/v1.8.31...v1.8.32
[1.8.31]: https://gitee.com/quant1x/engine.git/compare/v1.8.30...v1.8.31
[1.8.30]: https://gitee.com/quant1x/engine.git/compare/v1.8.29...v1.8.30
[1.8.29]: https://gitee.com/quant1x/engine.git/compare/v1.8.28...v1.8.29
[1.8.28]: https://gitee.com/quant1x/engine.git/compare/v1.8.27...v1.8.28
[1.8.27]: https://gitee.com/quant1x/engine.git/compare/v1.8.26...v1.8.27
[1.8.26]: https://gitee.com/quant1x/engine.git/compare/v1.8.25...v1.8.26
[1.8.25]: https://gitee.com/quant1x/engine.git/compare/v1.8.24...v1.8.25
[1.8.24]: https://gitee.com/quant1x/engine.git/compare/v1.8.23...v1.8.24
[1.8.23]: https://gitee.com/quant1x/engine.git/compare/v1.8.22...v1.8.23
[1.8.22]: https://gitee.com/quant1x/engine.git/compare/v1.8.21...v1.8.22
[1.8.21]: https://gitee.com/quant1x/engine.git/compare/v1.8.20...v1.8.21
[1.8.20]: https://gitee.com/quant1x/engine.git/compare/v1.8.19...v1.8.20
[1.8.19]: https://gitee.com/quant1x/engine.git/compare/v1.8.18...v1.8.19
[1.8.18]: https://gitee.com/quant1x/engine.git/compare/v1.8.17...v1.8.18
[1.8.17]: https://gitee.com/quant1x/engine.git/compare/v1.8.16...v1.8.17
[1.8.16]: https://gitee.com/quant1x/engine.git/compare/v1.8.15...v1.8.16
[1.8.15]: https://gitee.com/quant1x/engine.git/compare/v1.8.14...v1.8.15
[1.8.14]: https://gitee.com/quant1x/engine.git/compare/v1.8.13...v1.8.14
[1.8.13]: https://gitee.com/quant1x/engine.git/compare/v1.8.12...v1.8.13
[1.8.12]: https://gitee.com/quant1x/engine.git/compare/v1.8.11...v1.8.12
[1.8.11]: https://gitee.com/quant1x/engine.git/compare/v1.8.10...v1.8.11
[1.8.10]: https://gitee.com/quant1x/engine.git/compare/v1.8.9...v1.8.10
[1.8.9]: https://gitee.com/quant1x/engine.git/compare/v1.8.8...v1.8.9
[1.8.8]: https://gitee.com/quant1x/engine.git/compare/v1.8.7...v1.8.8
[1.8.7]: https://gitee.com/quant1x/engine.git/compare/v1.8.6...v1.8.7
[1.8.6]: https://gitee.com/quant1x/engine.git/compare/v1.8.5...v1.8.6
[1.8.5]: https://gitee.com/quant1x/engine.git/compare/v1.8.4...v1.8.5
[1.8.4]: https://gitee.com/quant1x/engine.git/compare/v1.8.3...v1.8.4
[1.8.3]: https://gitee.com/quant1x/engine.git/compare/v1.8.2...v1.8.3
[1.8.2]: https://gitee.com/quant1x/engine.git/compare/v1.8.1...v1.8.2
[1.8.1]: https://gitee.com/quant1x/engine.git/compare/v1.8.0...v1.8.1
[1.8.0]: https://gitee.com/quant1x/engine.git/compare/v1.7.9...v1.8.0
[1.7.9]: https://gitee.com/quant1x/engine.git/compare/v1.7.8...v1.7.9
[1.7.8]: https://gitee.com/quant1x/engine.git/compare/v1.7.7...v1.7.8
[1.7.7]: https://gitee.com/quant1x/engine.git/compare/v1.7.6...v1.7.7
[1.7.6]: https://gitee.com/quant1x/engine.git/compare/v1.7.5...v1.7.6
[1.7.5]: https://gitee.com/quant1x/engine.git/compare/v1.7.4...v1.7.5
[1.7.4]: https://gitee.com/quant1x/engine.git/compare/v1.7.3...v1.7.4
[1.7.3]: https://gitee.com/quant1x/engine.git/compare/v1.7.2...v1.7.3
[1.7.2]: https://gitee.com/quant1x/engine.git/compare/v1.7.1...v1.7.2
[1.7.1]: https://gitee.com/quant1x/engine.git/compare/v1.7.0...v1.7.1
[1.7.0]: https://gitee.com/quant1x/engine.git/compare/v1.6.9...v1.7.0
[1.6.9]: https://gitee.com/quant1x/engine.git/compare/v1.6.8...v1.6.9
[1.6.8]: https://gitee.com/quant1x/engine.git/compare/v1.6.7...v1.6.8
[1.6.7]: https://gitee.com/quant1x/engine.git/compare/v1.6.6...v1.6.7
[1.6.6]: https://gitee.com/quant1x/engine.git/compare/v1.6.5...v1.6.6
[1.6.5]: https://gitee.com/quant1x/engine.git/compare/v1.6.4...v1.6.5
[1.6.4]: https://gitee.com/quant1x/engine.git/compare/v1.6.3...v1.6.4
[1.6.3]: https://gitee.com/quant1x/engine.git/compare/v1.6.2...v1.6.3
[1.6.2]: https://gitee.com/quant1x/engine.git/compare/v1.6.1...v1.6.2
[1.6.1]: https://gitee.com/quant1x/engine.git/compare/v1.6.0...v1.6.1
[1.6.0]: https://gitee.com/quant1x/engine.git/compare/v1.5.9...v1.6.0
[1.5.9]: https://gitee.com/quant1x/engine.git/compare/v1.5.8...v1.5.9
[1.5.8]: https://gitee.com/quant1x/engine.git/compare/v1.5.7...v1.5.8
[1.5.7]: https://gitee.com/quant1x/engine.git/compare/v1.5.6...v1.5.7
[1.5.6]: https://gitee.com/quant1x/engine.git/compare/v1.5.5...v1.5.6
[1.5.5]: https://gitee.com/quant1x/engine.git/compare/v1.5.4...v1.5.5
[1.5.4]: https://gitee.com/quant1x/engine.git/compare/v1.5.3...v1.5.4
[1.5.3]: https://gitee.com/quant1x/engine.git/compare/v1.5.2...v1.5.3
[1.5.2]: https://gitee.com/quant1x/engine.git/compare/v1.5.1...v1.5.2
[1.5.1]: https://gitee.com/quant1x/engine.git/compare/v1.5.0...v1.5.1
[1.5.0]: https://gitee.com/quant1x/engine.git/compare/v1.4.9...v1.5.0
[1.4.9]: https://gitee.com/quant1x/engine.git/compare/v1.4.8...v1.4.9
[1.4.8]: https://gitee.com/quant1x/engine.git/compare/v1.4.7...v1.4.8
[1.4.7]: https://gitee.com/quant1x/engine.git/compare/v1.4.6...v1.4.7
[1.4.6]: https://gitee.com/quant1x/engine.git/compare/v1.4.5...v1.4.6
[1.4.5]: https://gitee.com/quant1x/engine.git/compare/v1.4.4...v1.4.5
[1.4.4]: https://gitee.com/quant1x/engine.git/compare/v1.4.3...v1.4.4
[1.4.3]: https://gitee.com/quant1x/engine.git/compare/v1.4.2...v1.4.3
[1.4.2]: https://gitee.com/quant1x/engine.git/compare/v1.4.1...v1.4.2
[1.4.1]: https://gitee.com/quant1x/engine.git/compare/v1.4.0...v1.4.1
[1.4.0]: https://gitee.com/quant1x/engine.git/compare/v1.3.9...v1.4.0
[1.3.9]: https://gitee.com/quant1x/engine.git/compare/v1.3.8...v1.3.9
[1.3.8]: https://gitee.com/quant1x/engine.git/compare/v1.3.7...v1.3.8
[1.3.7]: https://gitee.com/quant1x/engine.git/compare/v1.3.6...v1.3.7
[1.3.6]: https://gitee.com/quant1x/engine.git/compare/v1.3.5...v1.3.6
[1.3.5]: https://gitee.com/quant1x/engine.git/compare/v1.3.4...v1.3.5
[1.3.4]: https://gitee.com/quant1x/engine.git/compare/v1.3.3...v1.3.4
[1.3.3]: https://gitee.com/quant1x/engine.git/compare/v1.3.2...v1.3.3
[1.3.2]: https://gitee.com/quant1x/engine.git/compare/v1.3.1...v1.3.2
[1.3.1]: https://gitee.com/quant1x/engine.git/compare/v1.3.0...v1.3.1
[1.3.0]: https://gitee.com/quant1x/engine.git/compare/v1.2.9...v1.3.0
[1.2.9]: https://gitee.com/quant1x/engine.git/compare/v1.2.8...v1.2.9
[1.2.8]: https://gitee.com/quant1x/engine.git/compare/v1.2.7...v1.2.8
[1.2.7]: https://gitee.com/quant1x/engine.git/compare/v1.2.6...v1.2.7
[1.2.6]: https://gitee.com/quant1x/engine.git/compare/v1.2.5...v1.2.6
[1.2.5]: https://gitee.com/quant1x/engine.git/compare/v1.2.4...v1.2.5
[1.2.4]: https://gitee.com/quant1x/engine.git/compare/v1.2.3...v1.2.4
[1.2.3]: https://gitee.com/quant1x/engine.git/compare/v1.2.2...v1.2.3
[1.2.2]: https://gitee.com/quant1x/engine.git/compare/v1.2.1...v1.2.2
[1.2.1]: https://gitee.com/quant1x/engine.git/compare/v1.2.0...v1.2.1
[1.2.0]: https://gitee.com/quant1x/engine.git/compare/v1.1.9...v1.2.0
[1.1.9]: https://gitee.com/quant1x/engine.git/compare/v1.1.8...v1.1.9
[1.1.8]: https://gitee.com/quant1x/engine.git/compare/v1.1.7...v1.1.8
[1.1.7]: https://gitee.com/quant1x/engine.git/compare/v1.1.6...v1.1.7
[1.1.6]: https://gitee.com/quant1x/engine.git/compare/v1.1.5...v1.1.6
[1.1.5]: https://gitee.com/quant1x/engine.git/compare/v1.1.4...v1.1.5
[1.1.4]: https://gitee.com/quant1x/engine.git/compare/v1.1.3...v1.1.4
[1.1.3]: https://gitee.com/quant1x/engine.git/compare/v1.1.2...v1.1.3
[1.1.2]: https://gitee.com/quant1x/engine.git/compare/v1.1.1...v1.1.2
[1.1.1]: https://gitee.com/quant1x/engine.git/compare/v1.1.0...v1.1.1
[1.1.0]: https://gitee.com/quant1x/engine.git/compare/v1.0.9...v1.1.0
[1.0.9]: https://gitee.com/quant1x/engine.git/compare/v1.0.8...v1.0.9
[1.0.8]: https://gitee.com/quant1x/engine.git/compare/v1.0.7...v1.0.8
[1.0.7]: https://gitee.com/quant1x/engine.git/compare/v1.0.6...v1.0.7
[1.0.6]: https://gitee.com/quant1x/engine.git/compare/v1.0.5...v1.0.6
[1.0.5]: https://gitee.com/quant1x/engine.git/compare/v1.0.4...v1.0.5
[1.0.4]: https://gitee.com/quant1x/engine.git/compare/v1.0.3...v1.0.4
[1.0.3]: https://gitee.com/quant1x/engine.git/compare/v1.0.2...v1.0.3
[1.0.2]: https://gitee.com/quant1x/engine.git/compare/v1.0.1...v1.0.2
[1.0.1]: https://gitee.com/quant1x/engine.git/compare/v1.0.0...v1.0.1
[1.0.0]: https://gitee.com/quant1x/engine.git/compare/v0.9.9...v1.0.0
[0.9.9]: https://gitee.com/quant1x/engine.git/compare/v0.9.8...v0.9.9
[0.9.8]: https://gitee.com/quant1x/engine.git/compare/v0.9.7...v0.9.8
[0.9.7]: https://gitee.com/quant1x/engine.git/compare/v0.9.6...v0.9.7
[0.9.6]: https://gitee.com/quant1x/engine.git/compare/v0.9.5...v0.9.6
[0.9.5]: https://gitee.com/quant1x/engine.git/compare/v0.9.4...v0.9.5
[0.9.4]: https://gitee.com/quant1x/engine.git/compare/v0.9.3...v0.9.4
[0.9.3]: https://gitee.com/quant1x/engine.git/compare/v0.9.2...v0.9.3
[0.9.2]: https://gitee.com/quant1x/engine.git/compare/v0.9.1...v0.9.2
[0.9.1]: https://gitee.com/quant1x/engine.git/compare/v0.9.0...v0.9.1
[0.9.0]: https://gitee.com/quant1x/engine.git/compare/v0.8.9...v0.9.0
[0.8.9]: https://gitee.com/quant1x/engine.git/compare/v0.8.8...v0.8.9
[0.8.8]: https://gitee.com/quant1x/engine.git/compare/v0.8.7...v0.8.8
[0.8.7]: https://gitee.com/quant1x/engine.git/compare/v0.8.6...v0.8.7
[0.8.6]: https://gitee.com/quant1x/engine.git/compare/v0.8.5...v0.8.6
[0.8.5]: https://gitee.com/quant1x/engine.git/compare/v0.8.4...v0.8.5
[0.8.4]: https://gitee.com/quant1x/engine.git/compare/v0.8.3...v0.8.4
[0.8.3]: https://gitee.com/quant1x/engine.git/compare/v0.8.2...v0.8.3
[0.8.2]: https://gitee.com/quant1x/engine.git/compare/v0.8.1...v0.8.2
[0.8.1]: https://gitee.com/quant1x/engine.git/compare/v0.8.0...v0.8.1
[0.8.0]: https://gitee.com/quant1x/engine.git/compare/v0.7.9...v0.8.0
[0.7.9]: https://gitee.com/quant1x/engine.git/compare/v0.7.8...v0.7.9
[0.7.8]: https://gitee.com/quant1x/engine.git/compare/v0.7.7...v0.7.8
[0.7.7]: https://gitee.com/quant1x/engine.git/compare/v0.7.6...v0.7.7
[0.7.6]: https://gitee.com/quant1x/engine.git/compare/v0.7.5...v0.7.6
[0.7.5]: https://gitee.com/quant1x/engine.git/compare/v0.7.4...v0.7.5
[0.7.4]: https://gitee.com/quant1x/engine.git/compare/v0.7.3...v0.7.4
[0.7.3]: https://gitee.com/quant1x/engine.git/compare/v0.7.2...v0.7.3
[0.7.2]: https://gitee.com/quant1x/engine.git/compare/v0.7.1...v0.7.2
[0.7.1]: https://gitee.com/quant1x/engine.git/compare/v0.7.0...v0.7.1
[0.7.0]: https://gitee.com/quant1x/engine.git/compare/v0.6.9...v0.7.0
[0.6.9]: https://gitee.com/quant1x/engine.git/compare/v0.6.8...v0.6.9
[0.6.8]: https://gitee.com/quant1x/engine.git/compare/v0.6.7...v0.6.8
[0.6.7]: https://gitee.com/quant1x/engine.git/compare/v0.6.6...v0.6.7
[0.6.6]: https://gitee.com/quant1x/engine.git/compare/v0.6.5...v0.6.6
[0.6.5]: https://gitee.com/quant1x/engine.git/compare/v0.6.4...v0.6.5
[0.6.4]: https://gitee.com/quant1x/engine.git/compare/v0.6.3...v0.6.4
[0.6.3]: https://gitee.com/quant1x/engine.git/compare/v0.6.2...v0.6.3
[0.6.2]: https://gitee.com/quant1x/engine.git/compare/v0.6.1...v0.6.2
[0.6.1]: https://gitee.com/quant1x/engine.git/compare/v0.6.0...v0.6.1
[0.6.0]: https://gitee.com/quant1x/engine.git/compare/v0.5.9...v0.6.0
[0.5.9]: https://gitee.com/quant1x/engine.git/compare/v0.5.8...v0.5.9
[0.5.8]: https://gitee.com/quant1x/engine.git/compare/v0.5.7...v0.5.8
[0.5.7]: https://gitee.com/quant1x/engine.git/compare/v0.5.6...v0.5.7
[0.5.6]: https://gitee.com/quant1x/engine.git/compare/v0.5.5...v0.5.6
[0.5.5]: https://gitee.com/quant1x/engine.git/compare/v0.5.4...v0.5.5
[0.5.4]: https://gitee.com/quant1x/engine.git/compare/v0.5.3...v0.5.4
[0.5.3]: https://gitee.com/quant1x/engine.git/compare/v0.5.2...v0.5.3
[0.5.2]: https://gitee.com/quant1x/engine.git/compare/v0.5.1...v0.5.2
[0.5.1]: https://gitee.com/quant1x/engine.git/compare/v0.5.0...v0.5.1
[0.5.0]: https://gitee.com/quant1x/engine.git/compare/v0.4.9...v0.5.0
[0.4.9]: https://gitee.com/quant1x/engine.git/compare/v0.4.8...v0.4.9
[0.4.8]: https://gitee.com/quant1x/engine.git/compare/v0.4.7...v0.4.8
[0.4.7]: https://gitee.com/quant1x/engine.git/compare/v0.4.6...v0.4.7
[0.4.6]: https://gitee.com/quant1x/engine.git/compare/v0.4.5...v0.4.6
[0.4.5]: https://gitee.com/quant1x/engine.git/compare/v0.4.4...v0.4.5
[0.4.4]: https://gitee.com/quant1x/engine.git/compare/v0.4.3...v0.4.4
[0.4.3]: https://gitee.com/quant1x/engine.git/compare/v0.4.2...v0.4.3
[0.4.2]: https://gitee.com/quant1x/engine.git/compare/v0.4.1...v0.4.2
[0.4.1]: https://gitee.com/quant1x/engine.git/compare/v0.4.0...v0.4.1
[0.4.0]: https://gitee.com/quant1x/engine.git/compare/v0.3.9...v0.4.0
[0.3.9]: https://gitee.com/quant1x/engine.git/compare/v0.3.8...v0.3.9
[0.3.8]: https://gitee.com/quant1x/engine.git/compare/v0.3.7...v0.3.8
[0.3.7]: https://gitee.com/quant1x/engine.git/compare/v0.3.6...v0.3.7
[0.3.6]: https://gitee.com/quant1x/engine.git/compare/v0.3.5...v0.3.6
[0.3.5]: https://gitee.com/quant1x/engine.git/compare/v0.3.4...v0.3.5
[0.3.4]: https://gitee.com/quant1x/engine.git/compare/v0.3.3...v0.3.4
[0.3.3]: https://gitee.com/quant1x/engine.git/compare/v0.3.2...v0.3.3
[0.3.2]: https://gitee.com/quant1x/engine.git/compare/v0.3.1...v0.3.2
[0.3.1]: https://gitee.com/quant1x/engine.git/compare/v0.3.0...v0.3.1
[0.3.0]: https://gitee.com/quant1x/engine.git/compare/v0.2.9...v0.3.0
[0.2.9]: https://gitee.com/quant1x/engine.git/compare/v0.2.8...v0.2.9
[0.2.8]: https://gitee.com/quant1x/engine.git/compare/v0.2.7...v0.2.8
[0.2.7]: https://gitee.com/quant1x/engine.git/compare/v0.2.6...v0.2.7
[0.2.6]: https://gitee.com/quant1x/engine.git/compare/v0.2.5...v0.2.6
[0.2.5]: https://gitee.com/quant1x/engine.git/compare/v0.2.4...v0.2.5
[0.2.4]: https://gitee.com/quant1x/engine.git/compare/v0.2.3...v0.2.4
[0.2.3]: https://gitee.com/quant1x/engine.git/compare/v0.2.2...v0.2.3
[0.2.2]: https://gitee.com/quant1x/engine.git/compare/v0.2.1...v0.2.2
[0.2.1]: https://gitee.com/quant1x/engine.git/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/engine.git/compare/v0.1.9...v0.2.0
[0.1.9]: https://gitee.com/quant1x/engine.git/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/engine.git/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/engine.git/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/engine.git/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/engine.git/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/engine.git/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/engine.git/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/engine.git/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/engine.git/compare/v0.1.0...v0.1.1

[0.1.0]: https://gitee.com/quant1x/engine.git/releases/tag/v0.1.0

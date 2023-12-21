# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.9.1] - 2023-12-21
### Changed
- 调整规则配置, 增加NumberRange功能.
- 调整规则配置, 增加NumberRange功能.
- 修订订单备注确实的问题.

## [0.9.0] - 2023-12-20
### Changed
- 交易规则增加可止盈和可止损的判断方法.

## [0.8.9] - 2023-12-20
### Changed
- 调整ta-lib版本号.

## [0.8.8] - 2023-12-20
### Changed
- 统一数据到factors目录.
- 调整数据集源文件名.

## [0.8.7] - 2023-12-20
### Changed
- 调整117号卖出策略的逻辑.
- 调整交易规则中的止盈字段字段名.
- 调整box的字段注释.
- 交易规则增加止盈止损.
- 调整规则中市值的范围.
- 调整box部分字段注释.
- 清除废弃的代码.

## [0.8.6] - 2023-12-19
### Changed
- 开放exchange和box特征数据.
- 调整engine内置的特征组合文件名.

## [0.8.5] - 2023-12-19
### Changed
- 统一engine中关于行情快照结构体的定义.
- 整理部分代码.

## [0.8.4] - 2023-12-18
### Changed
- F10增加市值控制.
- 调整最小市值默认值为5亿.
- 修订测试策略的接口实现.

## [0.8.3] - 2023-12-18
### Changed
- History增加昨日开盘,收盘,成交量和成交额.

## [0.8.2] - 2023-12-18
### Changed
- 优化部分代码.
- 删除废弃的代码.
- 更新依赖库版本.
- 买入和卖出检查是否黑白名单.
- 删除废弃的代码.
- 简化一刀切卖出规则.
- 更新依赖库版本, gotdx的snapshot增加本地时间戳字段, 用以观察本地时钟和服务器的差距.
- 修订编译脚本, 从go.mod中读取module.
- 去除多余的echo.

## [0.8.1] - 2023-12-17
### Changed
- 调整获取最新数据的行数,明确数据类型.
- History增加9和19日的均价、均量线.
- 调整策略Strategy的评估接口, result参数改用泛型treemap.
- 增加使用go build -ldflags构建时传入的版本号的提示性注释.
- 新增命令行永久flag, avx2加速和cpu核数控制.
- 调整series函数.
- 增加交易时段测试代码.
- 性能分析默认调整为关闭状态.
- 特征数据增加验证样本的方法.
- 移除时间戳格式.
- 增加3个可能用到的时间戳格式.

## [0.8.0] - 2023-12-16
### Changed
- 新增黑白名单功能.
- 增加文档.

## [0.7.9] - 2023-12-15
### Changed
- 优化获取当前交易日期的方法.
- 服务方式运行更新快照去掉进度条.
- 屏蔽废弃的功能函数.
- 新增漏掉的卖出时段判断.
- 新增861号卖出策略编码ID的常量.
- 拟增加各Flag订单的总开关.
- 删除独立的卖出策略配置sell.
- 更新依赖库版本.
- 修订卖出策略中订单备注的信息.
- 修订卖出策略中订单备注的信息.
- 优化配置加载过程.
- 修订一刀切的常量名, 用Sell替代Sale.
- 修订配置文件的处理方法.
- 修订最新的配置项的默认配置文件.
- 拆分出runtime配置项, runtime参数增加debug开关, 默认关闭.
- 收敛engine中的所有command, 目的是为了统一处理debug开关.
- 调整卖出定时任务.
- 细化卖出策略.
- 增加增量计算均线的函数.
- 合并部分小功能代码.
- 删除废弃的代码.
- 更新依赖库版本.
- 新增价格笼子的计算方法.
- 新增拉取指定日期内策略首次买入的个股列表.
- 优化qmt策略名称的处理方法.
- 优化和调整交易规则结构.
- History增加2日均线和4日均线.
- 调整去重函数.
- 新增qmt相关的功能函数.
- 交易时段新增判断是否当天的最后一个交易时段.
- 交易规则增加持股周期字段.
- 新增通过策略编码获取用于qmt系统的string类型的策略名, 大写S作为前缀后面跟quant1x系统的策略编码.
- 修复code list可能存在变化的情况引发进度条异常的bug.
- 更新pandas版本号.

## [0.7.8] - 2023-12-13
### Changed
- 调整util工具包.
- 执行策略前输出策略概要.
- 收敛获取应用程序文件名的方法.
- 更新依赖库版本.
- 增加doc文档说明性源文件.
- 使用go 1.21.5.

## [0.7.7] - 2023-12-12
### Changed
- 实现一刀切的功能.
- 拆分下单委托函数, 支持直接填充策略名和订单备注.
- 更新依赖库版本.
- 新增计算涨跌停板价格的函数.
- 拟增加持仓列表功能.
- 调整测试代码.
- 策略方面, 增加一个特殊的卖出策略117(一刀切), 新增QMT体系中的策略名函数和订单备注函数.
- 调整订单状态机.
- 更新gotdx版本.
- 优化策略有效性判断.
- 按策略关联板块以及是否过滤两融.
- 修订TODO注释.
- 基础数据, 拟增加两融标的.
- 新增矫正策略交易时段的处理.
- 新增配置测试代码.
- 拆分出股票代码列表的函数.
- 交易员参数增加账户ID.
- 修订数据适配器相关的错误信息.
- 删除废弃的时间类测试代码.
- 调整定时任务的配置方式.
- 修订规则的错误信息.
- 调整策略的错误信息.
- 新增通过配置调整定时任务的开关和触发条件.
- 更新依赖库版本.
- 增加进入股票后直接向qmt proxy发起委托下单.
- 特征数据增加K线数据的最低要求限制的检查.
- 统一最低要求K线数量的常量为120.
- 增加账户和策略可用资金的计算方法.
- 优化编译脚本.
- 调整分割线.
- 增加mac和windows平台的amd64编译搅脚本.
- 优化买入卖出交易费用的计算方法.
- 增加通过预算输出交易费用对象.
- 拆分出交易配置对象.
- 新增交易费用结构体.
- 交易配置增加费率.
- 规则增加过滤股票代码前缀.
- 调整订单结构.
- 更新依赖库gox版本号.
- 调整交易参数.
- 交易员参数增加交易角色.
- 调整qmt持仓字段.
- 增加持仓结构.
- 增加qmt的常量.
- 订单字段增加tag式注释.
- F10增加财务数据报告期.
- 调整print指令, 增加输出缓存日期和特征日期.
- 更新依赖库pkg, ta-lib版本号.

## [0.7.6] - 2023-12-05
### Changed
- 新增一个只获取一只股票tick数据的函数, 目的是为了方便单元测试.
- 优化配置加载方式.
- 增加撤单时段.
- 更新依赖库gotdx,pkg版本号.
- 调整交易方向类型.
- 修复repair --all 特征数据不生效的bug.
- F10增加每股收益扣除字段.
- Proxy服务器地址.
- F10增加营业总收入.
- 修订日志中的错误描述.
- 修订日志中的错误描述.
- 增加委托和撤单两个函数.
- 增加查询委托功能.
- 删除废弃的交易参数结构体.
- 更新gox版本, http增加post方法支持.
- 调整http get请求的参数.
- 调整安全分http请求的方法.
- 修复错误的注释.
- 删除废弃的评估方法.
- 调整控制的最大订单数.
- 更新依赖库pandas版本号.

## [0.7.5] - 2023-12-03
### Changed
- 调整统计参数归于模型.
- 加载配置文件增加错误日志.
- 增加市场雷达功能.
- 增加交易参数配置.
- 调整目录结构.
- 增加流通市值.
- 从cache目录中拆分出config.
- 新增交易模块.
- 拟增加权限模块.

## [0.7.4] - 2023-11-27
### Changed
- 更新依赖库版本.
- 修复订单状态被覆盖的bug.

## [0.7.3] - 2023-11-26
### Changed
- 更新gox版本.
- 废弃部分字段.
- 删除废弃的测试代码.
- 剥离部分runtime功能到gox.
- 调整package.
- 调整部分快照相关的函数名.
- 统计结构体增加涨速字段.
- 股票池结构体增加活跃度和涨速字段.
- 增加1号策略通达信公式源代码.
- 配置文件新增定时任务开关.

## [0.7.2] - 2023-11-19
### Changed
- 优化规则分组.
- 删除独立的次新股规则.
- 实现strategy接口的1号策略, 去掉指针接收器的用法.
- 收敛长期不更新的依赖库到pkg.
- 拆分策略结果结构体为一个独立的源文件.
- 调整策略的执行方法, 改用注册的方式。暂时屏蔽回测功能。.
- 调整no1的方法顺序.
- 调整history csv字段名.
- 拟增加数据源切换功能.
- 增加公开函数的注释.
- 更新依赖库版本.
- 更新依赖库版本.
- 修订股票池.

## [0.7.1] - 2023-11-13
### Changed
- 更新gotdx版本号, 更新内置的板块数据文件.

## [0.7.0] - 2023-11-13
### Changed
- 修复振幅最小值的key重复的bug.

## [0.6.9] - 2023-11-13
### Changed
- 规则增加振幅范围0.00%~15.00%.

## [0.6.8] - 2023-11-13
### Changed
- 调整部分函数为公开.
- 增加一个预备的投票模块.
- 更新依赖库版本.

## [0.6.7] - 2023-11-07
### Changed
- 更新gotdx版本, 优化除权除息的股本变化处理方法.
- 修复股本变化的类型中遗漏送配股上市的bug.
- 更新其它非quant1x组织的依赖库版本.
- 更新ta-lib版本号.
- 调整基础过滤规则.
- 增加盘中实时订单标识.
- 增加无效周期的常量.
- 调整记分牌的tag.
- 增加更新内存中的K线操作.
- 增加定时任务配置结构体.
- 去除废弃的代码.

## [0.6.6] - 2023-10-31
### Changed
- 去除废弃的代码.
- 升级依赖库版本号.
- 调整规则和订单配置加载方式.
- 更新ta-lib版本.

## [0.6.5] - 2023-10-30
### Changed
- 调整MV5的计算方法.

## [0.6.4] - 2023-10-30
### Changed
- 调整精度条bar的序号处理方式.

## [0.6.3] - 2023-10-30
### Changed
- 数据2个空白行, 暂时先这么固定输出, 后面再优化调度任务结构.
- 调整bar的空白行.

## [0.6.2] - 2023-10-30
### Changed
- 增加业绩预告数据.

## [0.6.1] - 2023-10-29
### Changed
- 优化imports.
- 调整存储订单的同时输出到股票池.
- 调整股票池StockPool的package.
- 删除废弃的重置证券代码的处理方法.
- 收敛recover捕获panic异常的方法.
- 优化debug开关.
- 调整捕获panic的函数名.
- 增加Recover函数.
- 调整GoMaxProcs函数名.
- 调整GoMaxProcs函数路径.
- 调整证券名称、季报的处理方法.
- F10的公告信息只处理证券代码.
- History增加前5日分钟均量的方法.
- 优化1d缓存对象.
- 更新依赖库版本.
- 更新gotdx版本.
- 调整季报的处理方法.
- 增加从单个snapshot更新K线的函数.
- 新增具有滑动窗口速度控制的WaitGroup.
- 更新gox版本.
- 更新gotdx版本.
- 修复snapshot可能是nil的bug.
- 增加快照定时任务.
- 特征数据增加异常捕获.
- 删除废弃的pprof代码.
- 更新gox版本.
- 屏蔽实时更新K线的定时任务.
- 调度任务启动时增加互斥锁.
- 调整F10的证券名称字段.
- 次新股默认规则通过.
- 策略接口增加订单类型和过滤器.
- 增加策略文件缓存路径的常量.
- 数据集更新增加捕获异常.
- 修复map并发读写的bug.
- 分时数据增加异常捕获.
- 更新gotdx版本.
- 更新依赖库版本.
- 调整models.
- 增加leveldb测试代码.

## [0.6.0] - 2023-10-25
### Changed
- 增加调度任务日志.

## [0.5.9] - 2023-10-25
### Changed
- 调整F10的csv字段.

## [0.5.8] - 2023-10-25
### Changed
- 调整调度任务info级别日志内容.

## [0.5.7] - 2023-10-25
### Changed
- 优化实时更新K线的时间范围.

## [0.5.6] - 2023-10-25
### Changed
- 优化规则引擎.
- 子命令增加测试参数异常的测试性代码.
- 增加输出规则列表的子命令.
- ResourcesPath改为常量.
- 增加过滤规则功能.
- 主程序增加回测模块.
- 修订README中各模块的完成情况.
- 增加回测功能.

## [0.5.5] - 2023-10-24
### Changed
- 增加修订application的初始化代码.

## [0.5.4] - 2023-10-24
### Changed
- 完善宽表数据.
- 新增K线宽表数据的基础函数.

## [0.5.3] - 2023-10-24
### Changed
- 增加分时数据缓存.

## [0.5.2] - 2023-10-23
### Changed
- 修复xdxr缺少date和code的bug.

## [0.5.1] - 2023-10-23
### Changed
- 调整调度任务代码结构.
- 调整定时任务的回调函数.
- 调整服务接口.

## [0.5.0] - 2023-10-23
### Changed
- 调整实时更新K线的兜底逻辑.

## [0.4.9] - 2023-10-22
### Changed
- 调整定时任务的实现方式.
- 更新依赖库版本.
- 恢复非交易时段的不操作的逻辑.
- 实时更新K线增加内外盘两个字段.
- 增加实时更新K线.

## [0.4.8] - 2023-10-21
### Changed
- 调整等待应用结束的机制.
- 更新gox版本.

## [0.4.7] - 2023-10-21
### Changed
- 给pprof增加开关.
- Engine增加性能分析工具.
- 修订业绩预报结构体注释.
- 修正测试代码.

## [0.4.6] - 2023-10-20
### Changed
- 基础数据增加实时更新基础K线的函数.

## [0.4.5] - 2023-10-20
### Changed
- 增加系统服务子命令.
- 更新gox版本号.
- 清理废弃的代码.

## [0.4.4] - 2023-10-19
### Changed
- 增加daemon服务命令字.

## [0.4.3] - 2023-10-19
### Changed
- 更新依赖版本.

## [0.4.2] - 2023-10-19
### Changed
- 优化代码结构.

## [0.4.1] - 2023-10-19
### Changed
- 调整字段名.

## [0.4.0] - 2023-10-19
### Changed
- 调整数据接口归类划分.
- 调整数据接口归类划分.

## [0.3.9] - 2023-10-19
### Changed
- 调整数据接口归类划分.

## [0.3.8] - 2023-10-19
### Changed
- 细分数据接口.

## [0.3.7] - 2023-10-18
### Changed
- 调整csv字段名.
- 增加一个轻量的特性接口, 用来扩展子特征.

## [0.3.6] - 2023-10-18
### Changed
- 调整manifest结构体字段.

## [0.3.5] - 2023-10-18
### Changed
- 调整manifest结构体私有为公开.

## [0.3.4] - 2023-10-18
### Changed
- 提取抽象结构.
- 新增数据集和特征的manifest.

## [0.3.3] - 2023-10-17
### Changed
- 调整数据接口.
- 修订上一个季报没公布导致前十大流通股东列表为空的bug, 如果未公布, 应该沿用再上一个季度的数据.
- 优化代码.
- 修订cache1d结构体的注释.
- 收敛cache1d的缓存文件路径函数.
- 调整cache1d的new函数.
- 调整数据适配器接口的方法顺序.
- 删除废弃的代码.
- 调整数据适配器接口的方法顺序.
- 修订缓存适配器接口的注释.

## [0.3.2] - 2023-10-17
### Changed
- 调整源代码文件名.
- 调整代码结构.
- Aaa.
- 更新主要依赖库版本.
- 修订项目的主要关键词解释.
- 修订缓存操作接口的注释.
- 调整history结构的csv字段名.
- 调整HousNo1的csv字段名.
- 调整F10的csv字段名.
- 调整cache1d的缓存路径.
- 调整增量(不推荐)接口的package.
- 调整数据接口.
- 应用程序增加性能分析功能.
- 更新gox版本.
- 增加数据项接口.
- 增加数据运算接口.
- 增加忽略pprof文件.
- 调整记分牌的package.
- 调整提供者的方法名.
- 新增 数据的控制台命令支持接口.
- Trait特性接口增加提供者方法.
- 调整dataset方法.
- 调整dataset方法.
- 股票池增加策略状态字段.
- 新增规则接口.
- 新增summary和trait两个接口.
- 新增数据接口.
- 股票池增加规则字段.
- 增加股票池结构, 所有的数据都放在一个文件里面.

## [0.3.1] - 2023-10-13
### Changed
- 增加ants协程池控制并发数量.

## [0.3.0] - 2023-10-13
### Changed
- 测试协程方式跑特征数据.

## [0.2.9] - 2023-10-13
### Changed
- 优化update和repair数据处理流程.

## [0.2.8] - 2023-10-13
### Changed
- 优化update和repair数据处理流程.

## [0.2.7] - 2023-10-13
### Changed
- 优化update和repair数据处理流程.

## [0.2.6] - 2023-10-12
### Changed
- 增加周线,月线函数.

## [0.2.5] - 2023-10-12
### Changed
- 调整engine数据的提供者为engine.

## [0.2.4] - 2023-10-12
### Changed
- 命令字初始化改为显式.

## [0.2.3] - 2023-10-12
### Changed
- 调整子命令的检索逻辑.
- 新增数据验证check接口.
- 调整缓存的工厂用法.
- 新增通达信自选股列表导出函数.
- 增加一个单独的增量计算的接口备用.
- 增加通达信F10的资金流向, 这个数据因为网络请求的轮询机制, 数据很有可能存在不同源的问题, 从而导致数据不完整或者不正确.
- 变更源文件名.
- 调整历史成交记录的update和repair, 更新的日期应该采用cacheDate.
- 修改错误名.

## [0.2.2] - 2023-10-11
### Changed
- 更新gotdx版本, 历史成交数据去掉用pandas的方式读写, 改为切片和csv文件直接交换.

## [0.2.1] - 2023-10-11
### Changed
- 修订切片自动扩容地址变化引起的优先级较高的特征信息不能打印的bug.

## [0.2.0] - 2023-10-11
### Changed
- 子命令print自动检测是否打印特征数据, 暂时不支持结构嵌套.
- 修订bitmap结构体注释.
- 屏蔽暂时废弃的变量声明.
- 调整源文件名.
- 调整进度条的index.
- 基础数据增加历史成交数据.
- 数据插件增加get接口.
- 增加位图, 为将来扩展特征类型做准备.

## [0.1.9] - 2023-10-10
### Changed
- 将内部函数公开.

## [0.1.8] - 2023-10-10
### Changed
- 调整更新和修复子命令.
- 更新gox版本.
- 消除没有使用参数的告警提示.
- 标注废弃部分函数.
- 增加注释.

## [0.1.7] - 2023-10-10
### Changed
- 调整基本面数据的优先级.
- 修订项目总名称.
- 修订分支的描述.
- 调整变量的写法.
- 调整变量的写法.
- 调整变量的写法.
- 调整变量的写法.
- 调整变量的写法.
- 调整插件模式的遍历方法.
- 修订README, 增加对于协同开发方面的说明.
- 收录github.com/mattn/go-runewidth@v0.0.15.
- 移除测试性代码.
- 调整插件接口名.
- 增加插件接口, 用以收盘写数据操作.
- 增加smart接口.
- 修正安全分单词.
- 删除废弃的测试代码.
- 更新F10中公告的增持和减持的字段名.

## [0.1.6] - 2023-10-08
### Changed
- 删除废弃的特征组合box.

## [0.1.5] - 2023-10-08
### Changed
- 调整缓存机制的时间函数的package归属.
- 调整测试代码.
- 增加version, print子命令.
- 优化命令行参数解析.
- 更新依赖库的版本.
- 调整数据集合, 增加基础K线, 财报, 安全分, 除权除息.
- 新增东方财富数据的接口.
- 调整除权除息列表的测试代码.
- 增加通达信协议日期转换函数.
- 特征增加侯总1号策略.
- Repair增加特征数据.
- Repair增加基础数据.
- 增加异常是显示调用栈.
- 新增F10基本面特征数据组合.
- 增加个股安全评估数据.
- 修正cache1d的缓存关键字.
- 增补规范的文件名函数.
- 调整代码归属.
- 更新gox库版本.

## [0.1.4] - 2023-10-07
### Changed
- 更新gox、gotdx库版本.

## [0.1.3] - 2023-10-06
### Changed
- 调整数据集和特征组合.
- 执行策略之前增加同步即时行情数据的过程, 以便策略可以使用增量计算方法.
- 调整策略结果结构体字段顺序.
- 调整策略结果结构体.
- 更新gox版本.
- 拆分dataset.
- 调整基础数据集合.

## [0.1.2] - 2023-10-02
### Changed
- 完成第一个策略演示.

## [0.1.1] - 2023-10-01
### Changed
- 增加第一个策略执行的demo.
- Add ChangeLog.
- 增加趋势反转代码.
- 新增K线和除权除息的基础数据.

## [0.1.0] - 2023-09-27

### Changed

- 新增测试特征接口的代码, 以日K线为样本.
- 修订README.
- History增加日期的描述.
- 新增快照数据结构.
- 新增历史数据结构.
- 新增基础k线测试程序.
- 增加统一的常量模块.
- Add LICENSE.
- First commit.

[Unreleased]: https://gitee.com/quant1x/engine/compare/v0.9.1...HEAD

[0.9.1]: https://gitee.com/quant1x/engine/compare/v0.9.0...v0.9.1
[0.9.0]: https://gitee.com/quant1x/engine/compare/v0.8.9...v0.9.0
[0.8.9]: https://gitee.com/quant1x/engine/compare/v0.8.8...v0.8.9
[0.8.8]: https://gitee.com/quant1x/engine/compare/v0.8.7...v0.8.8
[0.8.7]: https://gitee.com/quant1x/engine/compare/v0.8.6...v0.8.7
[0.8.6]: https://gitee.com/quant1x/engine/compare/v0.8.5...v0.8.6
[0.8.5]: https://gitee.com/quant1x/engine/compare/v0.8.4...v0.8.5
[0.8.4]: https://gitee.com/quant1x/engine/compare/v0.8.3...v0.8.4
[0.8.3]: https://gitee.com/quant1x/engine/compare/v0.8.2...v0.8.3
[0.8.2]: https://gitee.com/quant1x/engine/compare/v0.8.1...v0.8.2
[0.8.1]: https://gitee.com/quant1x/engine/compare/v0.8.0...v0.8.1
[0.8.0]: https://gitee.com/quant1x/engine/compare/v0.7.9...v0.8.0
[0.7.9]: https://gitee.com/quant1x/engine/compare/v0.7.8...v0.7.9
[0.7.8]: https://gitee.com/quant1x/engine/compare/v0.7.7...v0.7.8
[0.7.7]: https://gitee.com/quant1x/engine/compare/v0.7.6...v0.7.7
[0.7.6]: https://gitee.com/quant1x/engine/compare/v0.7.5...v0.7.6
[0.7.5]: https://gitee.com/quant1x/engine/compare/v0.7.4...v0.7.5
[0.7.4]: https://gitee.com/quant1x/engine/compare/v0.7.3...v0.7.4
[0.7.3]: https://gitee.com/quant1x/engine/compare/v0.7.2...v0.7.3
[0.7.2]: https://gitee.com/quant1x/engine/compare/v0.7.1...v0.7.2
[0.7.1]: https://gitee.com/quant1x/engine/compare/v0.7.0...v0.7.1
[0.7.0]: https://gitee.com/quant1x/engine/compare/v0.6.9...v0.7.0
[0.6.9]: https://gitee.com/quant1x/engine/compare/v0.6.8...v0.6.9
[0.6.8]: https://gitee.com/quant1x/engine/compare/v0.6.7...v0.6.8
[0.6.7]: https://gitee.com/quant1x/engine/compare/v0.6.6...v0.6.7
[0.6.6]: https://gitee.com/quant1x/engine/compare/v0.6.5...v0.6.6
[0.6.5]: https://gitee.com/quant1x/engine/compare/v0.6.4...v0.6.5
[0.6.4]: https://gitee.com/quant1x/engine/compare/v0.6.3...v0.6.4
[0.6.3]: https://gitee.com/quant1x/engine/compare/v0.6.2...v0.6.3
[0.6.2]: https://gitee.com/quant1x/engine/compare/v0.6.1...v0.6.2
[0.6.1]: https://gitee.com/quant1x/engine/compare/v0.6.0...v0.6.1
[0.6.0]: https://gitee.com/quant1x/engine/compare/v0.5.9...v0.6.0
[0.5.9]: https://gitee.com/quant1x/engine/compare/v0.5.8...v0.5.9
[0.5.8]: https://gitee.com/quant1x/engine/compare/v0.5.7...v0.5.8
[0.5.7]: https://gitee.com/quant1x/engine/compare/v0.5.6...v0.5.7
[0.5.6]: https://gitee.com/quant1x/engine/compare/v0.5.5...v0.5.6
[0.5.5]: https://gitee.com/quant1x/engine/compare/v0.5.4...v0.5.5
[0.5.4]: https://gitee.com/quant1x/engine/compare/v0.5.3...v0.5.4
[0.5.3]: https://gitee.com/quant1x/engine/compare/v0.5.2...v0.5.3
[0.5.2]: https://gitee.com/quant1x/engine/compare/v0.5.1...v0.5.2
[0.5.1]: https://gitee.com/quant1x/engine/compare/v0.5.0...v0.5.1
[0.5.0]: https://gitee.com/quant1x/engine/compare/v0.4.9...v0.5.0
[0.4.9]: https://gitee.com/quant1x/engine/compare/v0.4.8...v0.4.9
[0.4.8]: https://gitee.com/quant1x/engine/compare/v0.4.7...v0.4.8
[0.4.7]: https://gitee.com/quant1x/engine/compare/v0.4.6...v0.4.7
[0.4.6]: https://gitee.com/quant1x/engine/compare/v0.4.5...v0.4.6
[0.4.5]: https://gitee.com/quant1x/engine/compare/v0.4.4...v0.4.5
[0.4.4]: https://gitee.com/quant1x/engine/compare/v0.4.3...v0.4.4
[0.4.3]: https://gitee.com/quant1x/engine/compare/v0.4.2...v0.4.3
[0.4.2]: https://gitee.com/quant1x/engine/compare/v0.4.1...v0.4.2
[0.4.1]: https://gitee.com/quant1x/engine/compare/v0.4.0...v0.4.1
[0.4.0]: https://gitee.com/quant1x/engine/compare/v0.3.9...v0.4.0
[0.3.9]: https://gitee.com/quant1x/engine/compare/v0.3.8...v0.3.9
[0.3.8]: https://gitee.com/quant1x/engine/compare/v0.3.7...v0.3.8
[0.3.7]: https://gitee.com/quant1x/engine/compare/v0.3.6...v0.3.7
[0.3.6]: https://gitee.com/quant1x/engine/compare/v0.3.5...v0.3.6
[0.3.5]: https://gitee.com/quant1x/engine/compare/v0.3.4...v0.3.5
[0.3.4]: https://gitee.com/quant1x/engine/compare/v0.3.3...v0.3.4
[0.3.3]: https://gitee.com/quant1x/engine/compare/v0.3.2...v0.3.3
[0.3.2]: https://gitee.com/quant1x/engine/compare/v0.3.1...v0.3.2
[0.3.1]: https://gitee.com/quant1x/engine/compare/v0.3.0...v0.3.1
[0.3.0]: https://gitee.com/quant1x/engine/compare/v0.2.9...v0.3.0
[0.2.9]: https://gitee.com/quant1x/engine/compare/v0.2.8...v0.2.9
[0.2.8]: https://gitee.com/quant1x/engine/compare/v0.2.7...v0.2.8
[0.2.7]: https://gitee.com/quant1x/engine/compare/v0.2.6...v0.2.7
[0.2.6]: https://gitee.com/quant1x/engine/compare/v0.2.5...v0.2.6
[0.2.5]: https://gitee.com/quant1x/engine/compare/v0.2.4...v0.2.5
[0.2.4]: https://gitee.com/quant1x/engine/compare/v0.2.3...v0.2.4
[0.2.3]: https://gitee.com/quant1x/engine/compare/v0.2.2...v0.2.3
[0.2.2]: https://gitee.com/quant1x/engine/compare/v0.2.1...v0.2.2
[0.2.1]: https://gitee.com/quant1x/engine/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/engine/compare/v0.1.9...v0.2.0
[0.1.9]: https://gitee.com/quant1x/engine/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/engine/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/engine/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/engine/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/engine/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/engine/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/engine/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/engine/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/engine/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/engine/releases/tag/v0.1.0

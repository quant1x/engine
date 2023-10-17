# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

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

[Unreleased]: https://gitee.com/quant1x/engine/compare/v0.3.2...HEAD

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

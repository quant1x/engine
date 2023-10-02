Quant1x-Engine
===

量化交易（数据类）引擎

## 1. 设计原则
- 更新历史数据，盘后更新历史数据
- 回补历史数据，对于新增特征组合或者因子要能回补历史数据，以便回测
- 盘中更新数据，盘中决策很重要，特征组合要有根据5档行情即时数据进行增量计算的能力
- 缓存数据必须具备按照日期切换的功能

## 2. 模块划分

| 级别 | 模块            | 功能                                                 | 盘中更新数据 | 更新当日数据 | 回补历史数据 |
|:---|:--------------|:---------------------------------------------------|:-------|:-------|:-----|
| 0  | cache         | 数据缓存                                               | [ ]    | [ ]    | [ ]  |
| 0  | factors       | 量化因子                                               | [ ]    | [ ]    | [ ]  |
| 0  | features      | 特征                                                 | [ ]    | [ ]    | [ ]  |
| 0 | tracker | 回测 | [ ]    | [ ]    | [ ] |

## 3. 使用示例

### 3.1 更新数据

```shell
engine update --all
```

### 3.2 补登历史特征数据

```shell
engine repair --history --start=20230101
```

### 3.3 执行1号策略

```shell
engine --strategy=1
```
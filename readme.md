## 立项本意与观点
* 历经多个大型系统开发项目，并逆向分析当前各种复杂度的站点、系统架构及其后端支撑团队组织结构认为大型服务系统更适合于采用面向过程与微服务结合架构进行构建系统，并组建对应支持团队，可有效降低项目复杂度及开发维护支出。
* 任何工具、系统、框架、编程语言应当适时收敛于特定稳定形态直至终结，不得任意无限延伸扩展。
* 期望构建一套方法使得系统开发、维护复杂度、难度随这套方法的升级不断降低，直至任何人都可以直接上手使用。
* 成为一套完善的网络系统开发框架、及附加支撑系统。
* 系统测试、发布运行、后期运维应当是一个统一的整体不分彼此。

## 弱水三千、只取一瓢
* 待更新

## 设计方法概述
* 面向过程定义通信、接口、存储
* 定义输入、输出标准描述结构与处理模式
* 结合测试、系统运行、运维定义系统运行模式
* 采用四层结构定义子系统及最终开放服务系统

## 四层系统结构说明如下
* 第一层：base 定义输入输出及机遇任务系统运行模式，附加基础定义和方法
* 第二层：rbox 定义执行逻辑与输入输出结果流转处理
* 第三层：fbox 真正的执行体，对输入与输出的处理
* 第四层：csbox 分布式、调度、存储与缓存、一致性、处理规则、可视化操作问题处理
* 当前base、rbox、fbox、csbox仍在验证开发中还会进行大的调整，当前所有模块均未开始优化
* 其中第一至三层构成单独进程模块处理单一任务；第四层为模块控制层并可在此层对下级模块再次组合封装
* 其中每一层都可更换细节，不变项为基础模式和定义；如：rbox、fbox（函数即服务的实现）均可根据业务目标进行重构


## 使用本套框架构造系统方法
* 待更新

## 演示系统构造与介绍
* 待更新

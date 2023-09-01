代码中有两种坐标，一种命名为 x y ，表示光标在终端窗口的位置坐标；
另一种命名为 row column(col) ，代表光标在文本中的位置。
一般来说， row 和 y 是相同的， column(col) 和 x 则不一定相同。
因为非 ASCII 字符（比如中文）会占用 2 个或者多个 x 。
两个坐标都是从 0 开始计数。

InputStream 负责解析输入，触发相应事件，调用 EventHandler 处理事件

EventHandler 负责处理事件，调用 Line 实现各种操作

Line 封装了文本和光标的操作

Document 保存文本和光标位置，提供一些辅助方法

Prompt 负责提示符的文本和样式

Screen 负责缓冲输出文本和样式，保存 xy 坐标维度的数据

Render 负责输出提示符和文本到屏幕

CommandLine 负责读取用户输入，串联整个流程。

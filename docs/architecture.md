# 坐标

代码中有两种坐标，一种命名为 x y ，表示光标在终端窗口的位置坐标；
另一种命名为 row column(col) ，代表光标在文本中的位置。
一般来说， row 和 y 是相同的， column(col) 和 x 则不一定相同。
因为非 ASCII 字符（比如中文）会占用 2 个或者多个 x 。
两个坐标都是从 0 开始计数。
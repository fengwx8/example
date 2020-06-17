# example

&emsp;&emsp;Exporter在原有的两个指标基础上增加了两个指标，metrics类型均为Gauge，分别是bitrate和server_state。
&emsp;&emsp;第一个bitrate命名不准确，我想用这个指标去模拟衡量服务端的当前链路可用带宽与链路总带宽的一个比值，以百分比数值表示。考虑到影响该数值最主要的因素应该是时间，所以直接用时间来模拟，不过这里只用了一个简单的抛物线做模拟，抛物线的自变量为小时，需要将60进制的分钟转换为10进制的小时，然后带入函数计算，就能得到对应的数值。
&emsp;&emsp;第二个指标server_state就相对简单了，因为以第一个为标准度量，由于bitrate是百分比数值，而且server_state表示服务器繁忙程度，所以直接用100减去第一个指标的值，就能得到繁忙程度的百分比数值表示。

# ProxifierForLinux

[English](README.md)

## 介绍

这个工具是一个利用redsocks和iptables实现类似于Proxifier的透明代理的工具。

## 用法

1.安装redsocks(https://github.com/darkk/redsocks)。
2. 添加 `redsocks` 组和 `redsocks` 用户：

````bash
groupadd redsocks
useradd -g redsocks redsocks
````

3. 以 root 身份运行 ProxifierForLinux。

````
sudo ./ProxifierForLinux
````

4. 添加规则并启动代理。

## 屏幕截图

![image-20241224085624035](./assets/image-20241224085624035.png)



![image-20241224085647392](./assets/image-20241224085647392.png)

![image-20241224085716495](./assets/image-20241224085716495.png)

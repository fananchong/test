# prometheus_test

## 安装部署

例子： [https://github.com/fananchong/test/blob/master/prometheus_test/setup.sh](https://github.com/fananchong/test/blob/master/prometheus_test/setup.sh)


## 自定义 Exporter && 采集指标

例子： [https://github.com/fananchong/test/blob/master/prometheus_test/metrics/main.go](https://github.com/fananchong/test/blob/master/prometheus_test/metrics/main.go)

## 告警

webhook 方式例子： [https://github.com/fananchong/test/blob/master/prometheus_test/alert/main.go](https://github.com/fananchong/test/blob/master/prometheus_test/alert/main.go)


## 可视化 Dashboard

Grafana 社区鼓励用户分享 Dashboard 通过 [https://grafana.com/dashboards](https://grafana.com/dashboards) 网站，可以找到大量可直接使用的 Dashboard

Grafana 中所有的 Dashboard 通过 JSON 进行共享，下载并且导入这些 JSON 文件，就可以直接使用这些已经定义好的 Dashboard

比如 [node-exporter_rev17.json](node-exporter_rev17.json) 就是从 [https://grafana.com/grafana/dashboards/8919](https://grafana.com/grafana/dashboards/8919) 页面上下载的


## 学习资料

- [https://yunlzheng.gitbook.io/prometheus-book](https://yunlzheng.gitbook.io/prometheus-book)

这里的一些练习，都是照着 prometheus-book 完成的，还有一些操作没文本方式保留下来，包括：
- 模板化 Dashboard
- prometheus 的服务发现功能

还有一些高级功能，仅了解下，未做练习，包括：
- prometheus 的集群部署与 HA 部署
- 基于 Kubernetes 使用 prometheus


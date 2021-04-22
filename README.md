# NVIDIA-Exporter

A costumed Prometheus exporter for process-level GPU utilization

## Usage

1. Build

    ```bash
    go build -o nvidia-exporter
    ```
2. Use (default listen_port is  `9114`)

	```bash
	./nvidia-exporter [listen_port]
	```

3. Test

   ```bash
   curl 127.0.0.1:[listen_port]/metrics
   ```

   ```yaml
   # HELP nvidia_mem_util mem_util
   # TYPE nvidia_mem_util gauge
   nvidia_mem_util{command="python",gpu_id="0",process_id="71767"} 6
   # HELP nvidia_sm_util sm_util
   # TYPE nvidia_sm_util gauge
   nvidia_sm_util{command="python",gpu_id="0",process_id="71767"} 20
   ```

## Link

- [nvidia_smi_exporter](https://github.com/zhebrak/nvidia_smi_exporter): Cluster-level GPU monitor with `nvidia-smi` 
- [prometheus-exporter](https://github.com/SongLee24/prometheus-exporter): Introduce basic concepts of Prometheus; Use Prometheus' Go client to implement an exporter 
timeout 180s bin/client -maddr=10.10.1.1 -writes=0.5 -c=10 -T=300
python3 client_metrics.py

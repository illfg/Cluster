go build

copy Cluster.exe C:\Users\Administrator\workspace\instances\instance1
copy Cluster.exe C:\Users\Administrator\workspace\instances\instance2
copy Cluster.exe C:\Users\Administrator\workspace\instances\instance3
copy Cluster.exe C:\Users\Administrator\workspace\instances\instance4
copy Cluster.exe C:\Users\Administrator\workspace\instances\instance5

start /d "C:\Users\Administrator\workspace\instances\instance1" Cluster.exe
start /d "C:\Users\Administrator\workspace\instances\instance2" Cluster.exe
start /d "C:\Users\Administrator\workspace\instances\instance3" Cluster.exe
start /d "C:\Users\Administrator\workspace\instances\instance4" Cluster.exe
start /d "C:\Users\Administrator\workspace\instances\instance5" Cluster.exe

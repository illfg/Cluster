go build

copy Cluster.exe C:\Users\Administrator\GolandProjects\testing\node1
copy Cluster.exe C:\Users\Administrator\GolandProjects\testing\node2
copy Cluster.exe C:\Users\Administrator\GolandProjects\testing\node3
copy Cluster.exe C:\Users\Administrator\GolandProjects\testing\node4

start /d "C:\Users\Administrator\GolandProjects\testing\node1" Cluster.exe
start /d "C:\Users\Administrator\GolandProjects\testing\node2" Cluster.exe
start /d "C:\Users\Administrator\GolandProjects\testing\node3" Cluster.exe
start /d "C:\Users\Administrator\GolandProjects\testing\node4" Cluster.exe
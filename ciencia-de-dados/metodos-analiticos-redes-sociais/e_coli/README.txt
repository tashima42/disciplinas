Data-Set: E.coli transcription regulation network
Description: Transcription regulation network in the E.coli 
Author: U. Alon & collaborators
Reference: S. Shen-Orr, R. Milo, S. Mangan  and U Alon, "Network
           motifs  in the transcriptional regulation network of
           Escherichia  coli.", Nature Genetics, 31:64-68 (2002). 
Additional-Info: This data set is part of the accompanying material of the book
                 "Complex Networks: Principles, Methods and Applications", 
                 V. Latora, V. Nicosia, G. Russo, Cambridge University Press (2017)
                 ISBN: 9781107103184
                 If you use the data set, please add a reference to
                 the original paper and to the book above.
Book-Reference: Chapter 8, Box 8.2

Filename: e_coli.net
Type: edge_list
Nodes: 424
Edges: 519
Directed: yes
Weighted: no
File-Format: "src_node dest_node type"
             where 'src_node' is the origin of the edge, 'dest_node'
             is the destination, and 'type' indicates the type of the
             interaction (1: activator, 2: repressor, 3: dual).
Comments: The original graph contained self-loops, but those have been
          removed from this data set.

Filename: e_coli_nodes.txt
Type: node_label_map
Rows: 424
File-Format: "node_label gene_name" 


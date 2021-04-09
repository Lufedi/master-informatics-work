#include <stdlib.h>
#include <stdio.h>
#include <omp.h>


int main(int argc, char **argv) {

    
    int m, n, q, c, d, k, sum, t = 0;
    int first[10][10], second[10][10], multiply[10][10];
 
  
    scanf("%d%d%d%d", &m, &n, &q, &t);
    
    for (c = 0; c < m; c++) {
        for (d = 0; d < n; d++) {
            scanf("%d", &first[c][d]);
        }
    }
        
    
    // elements of second matrix
    for (c = 0; c < n; c++)
      for (d = 0; d < q; d++)
        scanf("%d", &second[c][d]);
 
    for (c = 0; c < m; c++) {
      for (d = 0; d < q; d++) {
        for (k = 0; k < n; k++) {
          sum = sum + first[c][k]*second[k][d];
        }
 
        multiply[c][d] = sum;
        sum = 0;
      }
    
 
   }

    

    int numThreads = m*q;


#pragma omp parallel num_threads(numThreads)
{
    int i = omp_get_thread_num();
    int row = i / m;
    int col = i % q;

   
    
    int ss = 0;
    for (k = 0; k < n; k++) {
        ss = ss + first[row][k]*second[k][col];
        if (col == 0 && row == 0) { 
        }
    }

    multiply[row][col] = ss;
}


for (c = 0; c < m; c++) {
    for (d = 0; d < q; d++) {
        if (d == q - 1){
            printf("%d", multiply[c][d]);
        } else {
            printf("%d ", multiply[c][d]);
        }
    }
    printf("\n");
}

return 0;

}

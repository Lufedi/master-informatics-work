#include <stdlib.h>
#include <stdio.h>
#include <omp.h>
#include <time.h>

typedef int large;

const float EPSILPON = 0.000001;

large gcd(large fa, large fb) {
    large a = fa;
    large b = fb;
    large t;
    while (b != 0){
        t = b;
        b = a % b;
        a = t;
    }
    return a;
}

struct Fraction {
    large a;
    large b;
};

struct Fraction* newFraction(large a, large b, struct Fraction* result) {
    result->a = a;
    result->b = b;
}


struct Fraction MINUS_1 = {.a = 1, .b = 1};

void simplifyFraction(struct Fraction fraction) {
    large g = gcd(fraction.a, fraction.b);
    fraction.a = fraction.a / g;
    fraction.b = fraction.b / g;
}

int eq(struct Fraction f, struct Fraction o){
    float diff = abs(o.a/o.b - f.a/f.b);
    if(diff < EPSILPON) {
        return 1;
    }
    return 0;
}

void multiply(large f[2], large o[2], large res[2]) {
    res[0] = f[0]*o[0]; res[1]= f[1]*o[1];
}

void minus(large f[2], large o[2], large res[2]){
    res[0] = f[0]*o[1]-f[1]*o[0]; res[1] = f[1]*o[1];
}

void divide(large f[2], large o[2], large res[2]){
    res[0] = f[0]*o[1]; res[1] = f[1]*o[0];
}

void printMatrix(struct Fraction *matrix[10][10], int n) {
    int i,j;
    for(i = 0; i<n; i++) {
        for(j = 0 ; j < n+1; j++){
            printf("[%d,%d]", *matrix[i][j], *matrix[i][j]);
        }
        printf("\n");
    }
    printf("-------------");
}

large randNumber(){
    large upper = 2;
    large lower = 0;

    large num = (rand() %
           (upper - lower + 1)) + lower;
    return num;
}

void gauss(int N, int nthread){
    large a[N][N+1][2];
    large ratio[2];
    large x[N][2];
    int i, j, k;
    int f,o;

    // read data from stdin
    /*for(i=1; i<N; i++){
        for(j=1; j<=(N+1); j++) {
            scanf("%d", &f);
            scanf("%d", &o);
            struct Fraction fr = { .a=f, .b=o };
            a[i][j] = fr;
        }
    }*/
  

    for(i = 0; i < N; i++) {
        x[i][0] = 0; x[i][1] = 1;
        for(j = 0; j < N+1;j++) {
            a[i][j][0] = 0; a[i][j][1] = 1; 
        }
    }

   //FILE *file = fopen("/home/pipe/Documents/ECI/Master/PCYP/LearningGo/lab6a/input.txt", "r");
    for(i = 0; i < N; i++) {
        for(j = 0; j < N+1;j++) {
            //fscanf(file, "%d", fa);
            //fscanf(file, "%d", fb);
            a[i][j][0] = randNumber(); a[i][j][1] = 1; 
        }
    }

    large b[2];
    for(j=0; j<N; j++) {
        for(i=0; i<N; i++) {
            if(i!=j) {
                divide(a[i][j], a[j][j], b);
                #pragma omp parallel num_threads(nthread)
                {
                int segmentSize = N / nthread;
                int i = omp_get_thread_num();
                for(k=i*segmentSize; k < (i+1)*(segmentSize); k++) {
                    large temp[2], res[2];
                    multiply(b, a[j][k], temp);
                    minus(a[i][k],temp, res);
                    a[i][k][0] = res[0];
                    a[i][k][1] = res[1];
                }
                }
                /*for(k=1; k<=n+1; k++) { 
                    a[i][k] = minus(a[i][k], multiply(b, a[j][k]));
                }*/
            }
        }
    }
    for(i=0; i<N; i++) {
        large res[2];
        x[i][0] = 1;x[i][1] = 1;
        divide(a[i][N], a[i][i], res);
        x[i][0] = res[0]; x[i][1] = res[1];
    }
    /*
   //print
    for(i = 0; i<N; i++) {
        for(j = 0 ; j < N+1; j++){
            printf("[%d,%d]", a[i][j][0], a[i][j][1]);
        }
        printf("\n");
    }*/
 
}

static long get_nanos(void) {
    struct timespec ts;
    timespec_get(&ts, TIME_UTC);
    return (long)ts.tv_sec * 100000000000L + ts.tv_nsec;
}
int main() {
   
 //   fclose(file);
    int i,j,k;
    for(int i = 10; i  <  11; i++){
        int threads[] = {2,4,8,16, 64, 128, 256, 512, 1024, 2048};

        for(int j = 0; j < 10;j++) {
            double average = 0.0;

            for(int k = 0; k < 10; k++) {
                double time_spend = 0.0;
                clock_t begin = clock();
                gauss(i, threads[j]);

                clock_t end = clock();
                time_spend = (double)(end-begin) / CLOCKS_PER_SEC;
                average+=time_spend;
            }
            average= average / 10;
            printf("size %4d threads %4d  time %lf \n", i, threads[j] , average);
        }
    }
    return 0;
}
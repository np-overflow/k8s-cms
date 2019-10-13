#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <limits.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
    // allocate memory to simulate memory use
    size_t n_bytes = (size_t)510 * 1024 * 1024 * sizeof(char);
    char *mem = malloc(n_bytes);

    // perform numeric operations to load the cpu
    int x = 8;
    for(int i = INT_MIN; i < INT_MAX; i ++)
    {
        x = (x * x) / 2 % INT_MAX;
        x += 2345341;
        x = (x > 8) ? x: 8;
        mem[i % n_bytes] = (char)x;
    }
    
    printf("DONE\n");
    
    return 0;
}

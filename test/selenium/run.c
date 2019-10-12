#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <limits.h>

int main(int argc, char *argv[])
{  
    int x = 4;

    // perform numeric operations  forever to load the cpu
    for(int i = INT_MIN; i < INT_MAX; i ++)
    {
        x = (x * x) / 2 % INT_MAX;
    }
    
    printf("TEST END\n");
    return 0;
}

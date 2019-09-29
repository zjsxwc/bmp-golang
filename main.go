package main

/*
//
// Created by wangchao on 9/29/19.
//

#include <stdio.h>

//
//bmp_gen.h
//
#ifndef BMP_GEN_H
# define BMP_GEN_H

# include <stdio.h>
# include <stdlib.h>
# include <string.h>

# define SUCCESS 1
# define FAILURE 0

typedef	struct  	    s_color {
    double		        r;
    double		        g;
    double		        b;
}				        t_color;

#pragma pack(push, 1)
typedef struct	        s_bmp_file_header {
    unsigned char       bitmap_type[2];     // 2 bytes
    int                 file_size;          // 4 bytes
    short               reserved1;          // 2 bytes
    short               reserved2;          // 2 bytes
    unsigned int        offset_bits;        // 4 bytes
}				        t_bmp_file_header;
#pragma pack(pop)

#pragma pack(push, 1)
typedef struct	        s_bmp_image_header {
    unsigned int        size_header;        // 4 bytes
    unsigned int        width;              // 4 bytes
    unsigned int        height;             // 4 bytes
    short int           planes;             // 2 bytes
    short int           bit_count;          // 2 bytes
    unsigned int        compression;        // 4 bytes
    unsigned int        image_size;         // 4 bytes
    unsigned int        ppm_x;              // 4 bytes
    unsigned int        ppm_y;              // 4 bytes
    unsigned int        clr_used;           // 4 bytes
    unsigned int        clr_important;      // 4 bytes
}				        t_bmp_image_header;
#pragma pack(pop)

typedef struct	        s_bmp {
    FILE*               image;
    unsigned char*      data;
    t_bmp_file_header   fh;
    t_bmp_image_header  ih;
    int                 img_width;
    int                 img_height;
}				        t_bmp;

t_bmp*  bmpInit(int img_width, int img_height, const char *filename);
void    drawPixel(t_bmp *e, int x, int y, unsigned char r, unsigned char g, unsigned char b);
void    fileWrite(t_bmp *e);
void    clear(t_bmp *e);

#endif

//
//bmp_gen.c
//#include "bmp_gen.h"
//
void    bmpHeaderInit(t_bmp *e) {
    int         dpi = 96;
    int         image_size = e->img_width * e->img_height;
    int         file_size = 54 + 4 * image_size;
    int         ppm = dpi * 39.375;

    // Write File Header (14 bytes)
    memcpy(e->fh.bitmap_type, "BM", 2);
    e->fh.file_size = file_size;
    e->fh.reserved1 = 0;
    e->fh.reserved2 = 0;
    e->fh.offset_bits = 0x36;

    // Write Image Header (40 bytes)
    e->ih.size_header = sizeof(t_bmp_image_header);
    e->ih.width = e->img_width;
    e->ih.height = e->img_height;
    e->ih.planes = 1;
    e->ih.bit_count = 24;
    e->ih.compression = 0;
    e->ih.image_size = file_size;
    e->ih.ppm_x = ppm;
    e->ih.ppm_y = ppm;
    e->ih.clr_used = 0;
    e->ih.clr_important = 0;
}

// Open file and get a pointer to the FILE object created
int     fileWriteInit(t_bmp *e, const char *filename) {
    e->image = fopen(filename, "wb");
    if (e->image != NULL)
        return SUCCESS;
    else
        return FAILURE;
}

// Initialize structures (header / file) and memory allocation
t_bmp*  bmpInit(int img_width, int img_height, const char *filename) {
    t_bmp   *e;

    e = (t_bmp*)malloc(sizeof(t_bmp));
    e->data = (unsigned char*)malloc(sizeof(unsigned char) * (img_width * img_height * 3));
    if (e->data == NULL) {
        printf("Error during memory allocation !");
        return NULL;
    }
    e->img_height = img_height;
    e->img_width = img_width;
    bmpHeaderInit(e);

    if (fileWriteInit(e, filename) == FAILURE) {
        printf("Error during bitmap file creation !");
        return NULL;
    }
    return (e);
}


void    fileWrite(t_bmp *e) {
    unsigned char   bmppad[3] = {0,0,0};
    int             i;

    // Write headers into the file
    fwrite(&e->fh, 1, 14, e->image);
    fwrite(&e->ih, 1, 40, e->image);

    // Write data (+ padding when necessary)
    i = 0;
    for(i = 0; i < e->img_height; i++)
    {
        fwrite(e->data + (e->img_width * (e->img_height - i - 1) * 3), 3, e->img_width, e->image);
        fwrite(bmppad, 1, (4 - (e->img_width * 3) % 4) % 4, e->image);
    }
}

// Draw a pixel in a specific (x,y) position
void    drawPixel(t_bmp *e, int x, int y, unsigned char r, unsigned char g, unsigned char b) {
    int         i;
    int         j;

    i = x;
    j = (e->img_height - 1) - y;
    e->data[(i + j * e->img_height) * 3 + 2] = (unsigned char)(r);
    e->data[(i + j * e->img_height) * 3 + 1] = (unsigned char)(g);
    e->data[(i + j * e->img_height) * 3 + 0] = (unsigned char)(b);
}




typedef  struct s_memoryBmpResult {
    void* memoryHeadPtr;
    int size;
} t_memoryBmpResult;

// Initialize structures (header / file) and memory allocation
t_bmp*  bmpMemoryInit(int img_width, int img_height) {
    t_bmp   *e;

    e = (t_bmp*)malloc(sizeof(t_bmp));
    e->data = (unsigned char*)malloc(sizeof(unsigned char) * (img_width * img_height * 3));
    if (e->data == NULL) {
        printf("Error during memory allocation !");
        return NULL;
    }
    e->img_height = img_height;
    e->img_width = img_width;
    bmpHeaderInit(e);

    e->image = NULL;
    return (e);
}

t_memoryBmpResult memoryWrite(t_bmp *e) {
    unsigned char   bmppad[3] = {0,0,0};
    int             i;

    int totalSize = 0;
    totalSize += 14;
    totalSize += 40;
    i = 0;
    for(i = 0; i < e->img_height; i++)
    {
        totalSize += 3 * e->img_width;
        totalSize += (4 - (e->img_width * 3) % 4) % 4;
    }


    void* memoryPtr = malloc(totalSize);
    void* memoryHeadPtr = memoryPtr;
    memcpy(memoryPtr, &e->fh, 14);
    memoryPtr += 14;
    memcpy(memoryPtr, &e->ih, 40);
    memoryPtr += 40;
    i = 0;
    for(i = 0; i < e->img_height; i++)
    {
        memcpy(memoryPtr, e->data + (e->img_width * (e->img_height - i - 1) * 3), 3 * e->img_width);
        memoryPtr += 3 * e->img_width;

        memcpy(memoryPtr, bmppad, (4 - (e->img_width * 3) % 4) % 4);
        memoryPtr += (4 - (e->img_width * 3) % 4) % 4;
    }

    t_memoryBmpResult r;
    r.size = totalSize;
    r.memoryHeadPtr = memoryHeadPtr;
    return r;
}

void freeMemory(t_memoryBmpResult* rp)
{
    free(rp->memoryHeadPtr);
    rp->memoryHeadPtr = NULL;
    rp->size = -1;
}



//gcc testbmp.c
//./a.out
//open test.bmp

void testWriteFile() {
    printf("测试生成bmp图片文件\n");
    t_bmp *bp;
    int sideLength = 100;
    bp = bmpInit(sideLength,sideLength,"test.bmp");
    for (int i = 0; i < sideLength; ++i) {
        for (int j = 0; j < sideLength; ++j) {
            drawPixel(bp, i, j, 0x00, 0x00, 0x00);
        }
    }

    for (int i = 0; i < sideLength; ++i) {
        drawPixel(bp, i, i, 0xff, 0x00, 0x00);
    }
    for (int i = 0; i < sideLength; ++i) {
        drawPixel(bp, sideLength-i, i, 0xff, 0x00, 0x00);
    }
    fileWrite(bp);

    fclose(bp->image);

}


t_memoryBmpResult testMemory() {
    printf("测试生成内存bmp图片\n");
    t_bmp *bp;
    int sideLength = 100;
    bp = bmpMemoryInit(sideLength,sideLength);
    for (int i = 0; i < sideLength; ++i) {
        for (int j = 0; j < sideLength; ++j) {
            drawPixel(bp, i, j, 0x00, 0x00, 0x00);
        }
    }

    for (int i = 0; i < sideLength; ++i) {
        drawPixel(bp, i, i, 0xff, 0x00, 0x00);
    }
    for (int i = 0; i < sideLength; ++i) {
        drawPixel(bp, sideLength-i, i, 0xff, 0x00, 0x00);
    }
    t_memoryBmpResult r = memoryWrite(bp);

    FILE* testfile = fopen("test-from-c.bmp", "wb");
    fwrite(r.memoryHeadPtr, 1, r.size, testfile);
    fclose(testfile);

    return r;
}

int xmain() {//不能为main不然golang会报错说有两个main
    testWriteFile();
    testMemory();
    return 0;
}
 */
import "C"
import (
	"fmt"
)

func main() {
	r := C.testMemory()
	fmt.Println(r)
	C.freeMemory(&r)
	fmt.Println(r)
}

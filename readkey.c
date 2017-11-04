#include <termios.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdio.h>

void onKey(char *buf, int l);

struct termios oldterm;

void rk_init()
{
    // ungetc('\n', stdin);
    tcgetattr(STDIN_FILENO, &oldterm);
    struct termios newterm;
    newterm = oldterm;
    newterm.c_iflag &= ~(IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON);
    newterm.c_lflag &= ~(ECHO | ECHONL | ICANON | ISIG | IEXTEN);
    newterm.c_cflag &= (~(CSIZE | PARENB)) | CS8;
    newterm.c_cc[VMIN] = 1;
    newterm.c_cc[VTIME] = 0;
    tcsetattr(STDIN_FILENO, TCSANOW, &newterm);
    // printf("c\n");
}

void rk_release()
{
    tcsetattr(STDIN_FILENO, TCSANOW, &oldterm);
    // ungetc(ch, stdin);
}

int rk_read(char *buf, int len)
{
    return read(STDIN_FILENO, buf, len);
}

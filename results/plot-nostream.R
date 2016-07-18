all <- read.table('run5.clean', col.names=c('pos', 'total', 'val', 'type', 'size', 'stype'))
all$size <- factor(all$size)
all$val <- all$val/1000000
all$qps <- 1000/all$val

# The last bucket contains the total number of RPCs.
goodput <- all[seq(1,36)*101,]
goodput$goodput <- goodput$total*as.numeric(as.character(goodput$size))/1024/1024/60
# Consider zero payloads as being one bit.
goodput$goodput[goodput$size == 0] <-
    goodput$total[goodput$size == 0] / 8 / 1024/1024/60
goodput <- goodput[,c('type', 'size', 'stype', 'goodput')]
goodput <- goodput[goodput$stype == 'nostream',]

red <- c('#500003','#6C0D10','#87201C','#A0342A','#B74B3B','#CB624E','#DD7C64','#EA977C','#F3B397','#F7D1B4')
green <- c('#022609','#0E3C16','#205225','#356A36','#4D8148','#68995D','#85B174','#A6C98C','#C9E0A7','#EFF8C4')
sizeName <- function(val) {
  val <- as.numeric(val)
  if (val == 0) {
    return('0')
  }
  if (val < 1024*1024) {
    return(sprintf('%dK', val / 1024))
  }
  sprintf('%dM', val / 1024 / 1024)
}

plot.percent <- function(y, x1, x2, col) {
  rect(y-0.25, x1, y+0.25, x2, col=col, border=col)
}

plot.one <- function(typeArg, sizeArg, stypeArg, y, pallete) {
  d <- subset(all, type == typeArg & size == sizeArg & stype == stypeArg)
  d <- d$qps

  plot.percent(y, d[1], d[100], pallete[7])  # 98% of the data
  plot.percent(y, d[5], d[95],  pallete[5])  # 90% of the data
  plot.percent(y, d[25], d[75], pallete[1])  # 50% of the data
  segments(y-0.3, d[50], y+0.3, d[50], lwd=3, lend=1, col='white')
  mtext(sizeName(sizeArg), side=1, at=y, las=1, line=1)
  mtext(format(d[50], digits=2), side=1, at=y, las=1, line=2, cex=0.8)
  return(d[50])
}

png('plot-nostream.png', width=1200, height=800, bg='white')
par(family='DejaVu Sans')

# QPS plots.
par(mar=c(4, 4, 5, 0), fig=c(0, 0.5, 0.15, 1))

xrange <- c(0, 5000)
plot(c(0, 9), xrange, type='n', ann=F, axes=F, xaxs='i', yaxs='i')

i <- 1
for (size in levels(all$size)[1:4]) {
  m1 <- plot.one('v23', size, 'nostream', i, green)
  m2 <- plot.one('grpc', size, 'nostream', i+1, red)
  if (m1 > m2) {
    mtext(sprintf("+%.2fx", m1/m2), side=1, at=i, las=1, line=3, cex=0.8, col='gray')
    cat(sprintf("%7s %7.2f (+%4.2fx) %7.2f\n", sizeName(size), m1, m1/m2, m2))
  } else {
    mtext(sprintf("+%.2fx", m2/m1), side=1, at=i+1, las=1, line=3, cex=0.8, col='gray')
    cat(sprintf("%7s %7.2f          %7.2f (+%4.2fx)\n", sizeName(size), m1, m2, m2/m1))
  }
  i <- i + 2
}

x <- seq(xrange[1], xrange[2], 100)
axis(2, x, labels=F, las=2, lwd=0, lwd.tick=0.5, tcl=-0.3)
x <- seq(xrange[1], xrange[2], 1000)
axis(2, x, labels=x, lwd=0, las=1, lwd.tick=2)
mtext('Payload size \u2192', side=1, at=-0.15, las=1, line=1)
mtext('Median', side=1, at=-0.25, las=1, line=2, cex=0.8)
mtext('Improvement', side=1, at=-0.16, las=1, line=3, cex=0.8, col='gray')

par(mar=c(4, 4, 5, 2), fig=c(0.5, 1, 0.15, 1), new=T)
xrange <- c(0, 230)
plot(c(0, 11), xrange, type='n', ann=F, axes=F, xaxs='i', yaxs='i')

i <- 1
for (size in levels(all$size)[5:9]) {
  m1 <- plot.one('v23', size, 'nostream', i, green)
  m2 <- plot.one('grpc', size, 'nostream', i+1, red)
  if (m1 > m2) {
    mtext(sprintf("+%.2fx", m1/m2), side=1, at=i, las=1, line=3, cex=0.8, col='gray')
    cat(sprintf("%7s %7.2f (+%4.2fx) %7.2f\n", sizeName(size), m1, m1/m2, m2))
  } else {
    mtext(sprintf("+%.2fx", m2/m1), side=1, at=i+1, las=1, line=3, cex=0.8, col='gray')
    cat(sprintf("%7s %7.2f          %7.2f (+%4.2fx)\n", sizeName(size), m1, m2, m2/m1))
  }
  i <- i + 2
}

x <- seq(xrange[1], xrange[2], 2)
axis(2, x, labels=F, las=2, lwd=0, lwd.tick=0.5, tcl=-0.3)
x <- seq(xrange[1], xrange[2], 10)
axis(2, x, labels=format(x, nsmall=1), lwd=0, las=1, lwd.tick=2)

legend('topright', bty='n', cex=1.4,
       c('Vanadium 98%', 'Vanadium 90%', 'Vanadium 50%', 'Median',
         'gRPC 98%', 'gRPC 90%', 'gRPC 50%'),
       fill=c(green[7], green[5], green[1], 'white', red[7], red[5], red[1]))

# Goodput plots.
par(mar=c(2, 4, 0, 0), fig=c(0, 0.5, 0, 0.15), new=T)
xrange <- c(0, 260)
plot(c(0, 9), xrange, type='n', ann=F, axes=F, xaxs='i', yaxs='i')
x <- seq(1, 8)
rect(x-0.4, 0, x+0.4, goodput$goodput[x], col=c(green[7], red[7]), border=NA)
text(x, goodput$goodput[x], labels=sprintf("%0.2f", goodput$goodput[x]),
     pos=3, cex=0.8, col=c(green[2], red[2]))

xrange <- c(0, 200)
x <- seq(xrange[1], xrange[2], 10)
axis(2, x, labels=F, las=2, lwd=0, lwd.tick=0.5, tcl=-0.3)
x <- seq(xrange[1], xrange[2], 100)
axis(2, x, labels=format(x, nsmall=1), lwd=0, las=1, lwd.tick=2)

mtext('Improvement', side=1, at=-0.16, las=1, line=0.5, cex=0.8, col='gray')
x1 <- seq(1, 6, 2)
x2 <- seq(2, 6, 2)
axis(1, x2, labels=sprintf("+%.2fx", goodput$goodput[x2]/goodput$goodput[x1]),
     lwd=0, las=1, lwd.tick=0, line=-0.5, cex.axis=0.8, col.axis='gray')
x1 <- 7
x2 <- 8
axis(1, x1, labels=sprintf("+%.2fx", goodput$goodput[x1]/goodput$goodput[x2]),
     lwd=0, las=1, lwd.tick=0, line=-0.5, cex.axis=0.8, col.axis='gray')

legend('topleft', bty='n', inset=c(0, 0.1),
       c('Vanadium goodput (MiB)', 'gRPC goodput (MiB)'),
       pch=15, col=c(green[6], red[6]))

par(mar=c(2, 4, 0, 2), fig=c(0.5, 1, 0, 0.15), new=T)
xrange <- c(0, 260)
plot(c(0, 11), xrange, type='n', ann=F, axes=F, xaxs='i', yaxs='i')
x <- seq(1, 10)
rect(x-0.4, 0, x+0.4, goodput$goodput[x+8], col=c(green[7], red[7]), border=NA)
text(x, goodput$goodput[x+8], labels=sprintf("%0.2f", goodput$goodput[x+8]),
     pos=3, cex=0.8, col=c(green[2], red[2]))

xrange <- c(0, 200)
x <- seq(xrange[1], xrange[2], 10)
axis(2, x, labels=F, las=2, lwd=0, lwd.tick=0.5, tcl=-0.3)
x <- seq(xrange[1], xrange[2], 100)
axis(2, x, labels=format(x, nsmall=1), lwd=0, las=1, lwd.tick=2)

x1 <- seq(1, 10, 2)
x2 <- seq(2, 10, 2)
axis(1, x1, labels=sprintf("+%.2fx", goodput$goodput[x1+8]/goodput$goodput[x2+8]),
     lwd=0, las=1, lwd.tick=0, line=-0.5, cex.axis=0.8, col.axis='gray')

# Main title.
par(mar=c(6, 0, 3, 0), fig=c(0, 1, 0, 1), new=T)
title(main='Vanadium vs gRPC QPS for non-streaming RPCs', cex.main=2)
dev.off()

{89K趋势策略线, MA1X, V1.55, 2022-12-03}

{----------------< 基础数据部分 >----------------}
{比对周期长度}
C_N:=5;
{量价关系最低校对比率}
C_S:=0.191;
{涂画与股价的纵向比率}
C_PX:=0.002;
{真阳线}
C_ISMALE:=CLOSE > REF(CLOSE,1);
{成交量较上一个周期放大}
C_VOL:= VOL>REF(VOL,1);
{成交量均线周期}
VOL_PERIOD:=5;
{成交量比例}
VOLSCALE:=1+C_S;
{高股价或指数的计算方法, 比MAVOL5高出C_S/10且比前一日方法}
X_INDEX:=VOL>=MA(VOL, VOL_PERIOD)*(1 + C_S/10);
{一般股价的计算方法}
X_GENERAL:=VOL>=MA(VOL, VOL_PERIOD)*VOLSCALE;
{指数类或者高股价类的成交量不太可能像个股那样成倍放量, 这里做一个降级处理}
X:=IFF(CLOSE>=100, X_INDEX, X_GENERAL) AND C_ISMALE AND C_VOL;
{确定X上一次确立成立是在哪一天}
DN: BARSLAST(X),NODRAW;

{强弱提示}
XWIDTH:=0.009499432279;
PI:=3.141592654;
MA3:=MA(CLOSE,5);
QR1:=CLOSE-MA3;
QR2:=QR1/MA3;
QR3:=QR2/XWIDTH;
QR4:=ATAN(QR3);
QR:QR4*180/PI,NODRAW,COLORLIGRAY;
{----------------< A: MA移动平均线 >----------------}
P1:=5;
P2:=13;
P3:=21;
P4:=34;
P5:=65;
P6:=89;
P7:=144;
MA5:MA(CLOSE,P1),COLORFF8D1E;
MA13:MA(CLOSE,P2),COLOR0CAEE6;
MA21:MA(CLOSE,P3),COLORE970DC;
MA34:MA(CLOSE,P4),COLOR0080FF;
MA65:MA(CLOSE,P5),COLOR22C57E;
MA89:MA(CLOSE,P6), COLORCCFF00;
MA144:MA(CLOSE,P7),COLOR13FFFF;

{----------------< C: 止盈止损线 >----------------}
{放量上攻作为一个小阶段的起点, 该起点K线的最低价作为止盈止损线}
ZS_VOL:=REF(VOL,DN);
ZS_LOW:=REF(LOW,DN);
{1. 画止损线, 跌破止损线要清仓, 浅黄色虚线}
止损:DRAWLINE(VOL=ZS_VOL,ZS_LOW,DN>0,ZS_LOW,1), DOTLINE,COLOR13FFFF;
{止损: ZS_LOW,DOTLINE, LINETHICK1,COLORCCFF00;}
{X0:IFF(止损<0, -1, IFF(止损=0,0,1));}
{2. 输出最低价}
DRAWNUMBER(VOL=ZS_VOL AND DN=0, LOW*(1-C_PX), ZS_LOW), COLOR13FFFF;
{3. 无法显示'突破', 周期已被重置, 需要另想办法}
{DRAWTEXT(CLOSE>ZS_VOL AND DN>0, HIGH*(1+C_PX), '突破');}

{----------------< D: 压力线 >----------------}
{D: 压力线}
{0. 压力线, 备用}
{H0:=FINDHIGH(HIGH,0,5,1);
N31: FINDHIGHBARS(HIGH,0,5,1),NODRAW;
V2:=REF(VOL,N31);
N3:BARSLAST(VOL=V2 AND N31>0)+1,NODRAW;
H1:REF(HIGH,N3),NODRAW;
}
H0:=HHV(HIGH, C_N);
YL_N:BARSLAST(C_ISMALE AND HIGH=H0 AND C_VOL),NODRAW;
YL_VOL:=REF(VOL,YL_N);
YL_HIGH:=REF(HIGH,YL_N);
{1. 画压力线, 回踩后拉升收盘有突破的, 第2天买入, 浅蓝色虚线}
压力:DRAWLINE(VOL=YL_VOL,YL_HIGH,YL_N>0,YL_HIGH,1), DOTLINE, COLORCCFF00;
{2. 输出最高价}
DRAWNUMBER(VOL=YL_VOL AND YL_N=0, HIGH*(1+C_PX), YL_HIGH), COLORCCFF00;

{----------------< E: 箱体 >----------------}
ZN:=DN,NODRAW;
ZS:=ZS_LOW,NODRAW;
YN:=YL_N,NODRAW;
YH:=YL_HIGH,NODRAW;
YV:=YL_VOL,NODRAW;

{技术性测试支撑}
支撑周期:=IFF(ZN=0,REF(ZS,1)+1,ZN);
支撑线:=IFF(ZN=0,REF(ZS,1),ZS);
昨日最高:=REF(HIGH,1);
今日最低:=LOW;
洗盘力度:=100*(昨日最高-今日最低)/昨日最高;
回踩:支撑周期>=2 AND 洗盘力度>=2.33 AND LOW<支撑线 AND CLOSE>=支撑线,COLORYELLOW, NODRAW;
DRAWICON(回踩>0,LOW,1);

{技术性测试抛压}
压力周期:=IFF(YN=0,REF(YN,1)+1,YN),NODRAW;
压力线:=IFF(YN=0,REF(YH,1),YH);
压力量:=IFF(YN=0,REF(YV,1),YV);
T1:=REF(CLOSE,1)<压力线 OR REF(LOW,1)<压力线 OR LOW<压力线,NODRAW;
T2:=CLOSE>=压力线,NODRAW;
{T3:=VOL>测试压力量,NODRAW;}
突破:T1 AND T2 AND 压力周期>=2,COLORLIMAGENTA,NODRAW;
DRAWICON(突破>0,LOW,38);

{----------------< F: 线性回归 >----------------};
HG10:REF(YH,YN+1),NODRAW;
HG11:BARSLAST(HIGH=HG10*1),NODRAW;
HG20:BARSLAST(HIGH=YH),NODRAW;
{REF(YN,1),NODRAW;
HG02:YN,NODRAW;
HG03:REF(YN,1),NODRAW;
HG04:YH-REF(YH,1),NODRAW;
HG05:
HG1:SLOPE(YH,REF(YH,1));
HG2:CLOSE*HG1;}

{----------------< Z: 测试代码 >----------------};

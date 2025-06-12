# E 魔法首行：启动 HerCode 的温柔模式

#----------------------------------------
# E 人口区块 start
#   HerCode约定：程序从start：开始执行
#-----------------------------------------

#----------------------------------------
#「 你可以做到 (you_can_do_this)
#   这是一段 *鼓励式函数”:
#   - 功能：向终端打印两句话，欢迎女孩们来到专属的编程世界
#   - 关键字说明：
#       function   ----定义函数（像在写日记那样自然）
#       say        ---- 输出一句温暖的话
#       end        ---- 结束函数或代码块
#----------------------------------------

function you_can_do_this:
say "Hello! Her World!"      #终端打印友好的问候
say "编程很美，也属于你！"   #再来一句给自己的加油
end

function you_can_do_this1 food:
say "编程很美，也属于你！"      #再来一句给自己的加油
say "今天给自己加餐: " + food
end

function you_can_do_this2:
say "Hello! Her # World!"      #终端打印友好的问候
say "编程很美，也属于你！"   #再来一句给自己的加油
end

#----------------------------------------
#「 颜值评级函数
#   根据颜值评级
#----------------------------------------
function grade score:
    if score >= 90
        say "姐妹你太漂亮了! 颜值有 " + score
    else
        say "姐妹你好漂亮! 颜值有 " + score
    endif
end

#----------------------------------------
# E 人口区块 start
#   HerCode约定：程序从start：开始执行
#-----------------------------------------

start:
you_can_do_this                 #调用鼓励函数
you_can_do_this1("我吃柠檬")    #调用鼓励函数
you_can_do_this2                #调用鼓励函数

say "不吃香菜"                  #真不喜欢吃

var i = 1
var sum = 0
while i <= 10
    sum = sum + i
    i = i + 1
endif

say "1到10的和: " + sum # 对我的颜值进行评级

grade(255)              # 对我的颜值进行评级

grade(80)               # 对我的颜值进行评级

end
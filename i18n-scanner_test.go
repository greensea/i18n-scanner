package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	strs := Parse(sample, "__")
	if len(strs) != 12 {
		t.FailNow()
	}
}

func TestFile(t *testing.T) {
	f := NewFile()
	d := os.TempDir()
	p := fmt.Sprintf("%s/i18n-scanner_test-%d.json", d, time.Now().UnixNano())

	f.Load(p)
	f.AddLocale("zh")
	f.AddLocale("en")
	f.Add("Test Message")
	f.Add("测试翻译条目")
	f.Save(p)

	defer os.Remove(p)

	f.Load(p)
}

var sample string = `<script setup>
import {ref} from 'vue'
</script>
<script>
let speeds = []
for (let i =15; i < 32; i++) {
    speeds.push({
        hardness: 2**i - 1,
        hardness_html: '2<sup>${i}</sup> - 1  <span class="text-gray-400"> == ${2**i - 1}</span>',
        time_str: this.__("未计算"),
    })
}

export default {
    data() {
        return {
            activeName: "performance",
            speeds: speeds,
            isEvaling: {},
            n: "9e471153ecd402c1779e21985cca2ef2df1ff567bc1767f559256f3f592f4f7669d34e7b959d259796307040f8fc98c2702049e48730c9a4dcb2e4d66d60b5f3",
            x: "62bbf3bbbad2b415",
        }

    },

    mounted() {
        
    },

    methods: {
        async startEval(row) {
            let t = row.hardness
            let worker = new Worker("vdfworker.js")
            let that = this
            let stime = (new Date()).getTime()

            this.isEvaling[t] = true

            worker.onmessage = msg => {
                let d = msg.data;

                if (d[0] == "eval_progress") {
                    console.debug("Got VDF calculate progress", d[1])
                    row.time_str = Math.round(d[1] * 100) + "%"
                } else if (d[0] == "eval_result") {
                    row.time_ms = (new Date()).getTime() - stime
                    if (row.time_ms > 1000) {
                        row.time_str = (row.time_ms / 1000).toFixed(1) + that.__(" 秒")
                    } else {
                        row.time_str = Math.round(row.time_ms) + that.__(" 毫秒")
                    }
                    this.isEvaling[t] = false
                    worker.terminate()
                } else {
                    console.log(d);
                }
            }

            worker.postMessage(['eval', this.x, t, this.n])

//            worker.terminate()



        },
    },

    computed: {

    },
}
</script>


<template>

    <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick" tab-position="left">
        <el-tab-pane :label="__(    '原理说明')" name="intro" class="intro">
            
            <!-- <p class="text-center">$y \equiv x ^ {2 ^ t} (\bmod n)$</p> -->
            <p class="text-center">y = x<sup>2<sup>t</sup></sup> mod n</p>
            
        </el-tab-pane>

        <el-tab-pane :label="__('API 列表')" name="api" class="api" >
            <h2>wcaptcha(api_key)</h2>
            
            <highlightjs code='var w = new wcaptcha("YOUR_API_KEY");'></highlightjs>
            
            <highlightjs code='var w = new wcaptcha("YOUR_API_KEY");
w.bind("#any-dom-selector");'></highlightjs>
            <el-divider />

            <h2>wcaptcha.prototype.prove()</h2>
            <highlightjs code='var w = new wcaptcha("YOUR_API_KEY");
w.prove().then(proof => {
    console.log("proof is", prove);
}");'></highlightjs>

            <el-divider />
            <h2>wcaptcha.prototype.getProblem()</h2>
        </el-tab-pane>

        <el-tab-pane label="自定义使用" name="advance" class="advance">

            <highlightjs code="w = new wcaptcha(API_KEY);
w.onprogress = (progress) => {
    console.log('进度百分比:', progress * 100)
}"></highlightjs>

            <highlightjs code='w.prove().then(proof => {
    formData = new FormData()
    formData.append("wcaptcha-proof", proof)
    formData.append("Your", "Form Data")
    
    // Submit your data...
})
            '></highlightjs>

            

        </el-tab-pane>

        <el-tab-pane :label="__('速度测试')" name="performance" class="performance">

            <el-table :data="speeds" class="mt-2 max-w-md" table-layout="auto">
                <el-table-column :label="__('难度')">
                    <template #default="scope">
                        <span v-html="scope.row.hardness_html"></span>
                    </template>
                </el-table-column>
                <el-table-column prop="time_str" :label="__('验证耗时')" />
                <el-table-column prop="str" :label="__('测试')">
                    <template #default="scope">
                        <el-button @click="startEval(scope.row)" :loading="isEvaling[scope.row.hardness]" :disabled="isEvaling[scope.row.hardness]">
                            <span v-if="(isEvaling[scope.row.hardness] == true)">{{__('正在计算')}}</span>
                            <span v-else>{{__('开始测试')}}</span>
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>

        </el-tab-pane>

        <el-tab-pane :label="__('常见问题')" name="qna" class="qna">
 
        </el-tab-pane>

    </el-tabs>

</template>


<style scoped>
.intro p, el-ul, ol {
    line-height: 1.5rem;
    margin: 1em;
    text-indent: 2em;
}

.api h2, .qna h2 {
    font-size: 1.2rem;
    font-weight: bold;
    margin-bottom: 0.5em;
}
.api p, .qna p, .advance p, .performance p {
    margin-bottom: 0.5em;
}
</style>
`

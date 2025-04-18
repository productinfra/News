<template>
	<div>
		<el-card>
			<div slot="header" class="clearfix">
				<div class="web-status">🚀 Status</div>
			</div>
			<div class="run-time"></div>
		</el-card>
	</div>
</template>

<script>
import { runTime } from "@/utils/runTime"
export default {
	name: 'TimeMeter',
	data() {
		return {
			runTime: new Date("2025-4-8")
		}
	},
	mounted() {
		this.init();
	},
	methods: {
		init() {
			runTime({
				el: '.run-time',
				timeStamp: this.runTime.getTime(),
				desc: 'Up time: '
			});
		}
	}
}
</script>

<style lang="less" scoped>
.web-status {
	font-weight: 600;
	text-align: center;
}

.run-time {
	background-image: linear-gradient(to right, orange, purple);
	-webkit-background-clip: text;
	color: transparent;
	text-align: center;
}
</style>
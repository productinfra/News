<template>
	<div>
		<el-row>
			<el-col :span="24">
				<div class="grid-content bg-purple-dark footer">
					<ul class="friends-link">
						<li>
							<a href="https://github.com/productinfra/News">Github</a>
						</li>
					</ul>
					<h2 class="compony">Team 10</h2>
				</div>
			</el-col>
		</el-row>
	</div>
</template>

<script>
export default {
	name: "Footer",
};
</script>

<style>
* {
	list-style: none;
	text-decoration: none;
}

.footer {
	min-height: 100px;
	margin-top: 20px;
	border-top: 5px solid #74dfff;
	background-color: rgba(43, 42, 42, 1);
}

.friends-link {
	width: 50%;
	height: 50px;
	margin: 0 auto;
	line-height: 50px;
	display: flex;
	justify-content: space-around;
}

.friends-link li {
	width: 25%;
	height: 100%;
	color: #fff;
}

.friends-link li a {
	display: block;
	width: 100%;
	height: 100%;
	color: #fff;
	text-align: center;
}

.friends-link li a:hover {
	color: orange;
	text-decoration: underline;
}

.footer .compony {
	width: 50%;
	height: 50px;
	margin: 0 auto;
	color: #fff;
	line-height: 50px;
	text-align: center;
}</style>
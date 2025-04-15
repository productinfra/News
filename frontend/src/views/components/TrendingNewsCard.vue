<template>
	<el-row>
		<el-col :span="24">
			<el-card :body-style="{ padding: 0, border: 'none' }" class="box-card">
				<h2 class="title">{{ title }}</h2>
				<ul class="trending-news-list">
					<li class="trending-news-item" v-for="(news, index) in displayedList" :key="news.title">
						<div class="introduction">
							<span class="rank">#{{ index + 1 }}</span>
							<a :href="news.url" target="_blank" rel="noopener noreferrer">{{ news.title }}</a>
						</div>
						<el-divider></el-divider>
					</li>
					<div class="next-btn" v-if="newsList.length > 0" @click="getMorenews">
						<i :class="loading"></i>
						<span>Load more...</span>
					</div>
					<el-skeleton v-if="newsList.length <= 0" :rows="6" />
				</ul>
			</el-card>
		</el-col>
	</el-row>
</template>

<script>
export default {
	name: 'TrendingNewsCard',
	props: {
		title: {
			type: String,
			require: true,
			default: 'Trending News'
		},
		language: {
			type: String,
			require: true,
			default: 'all'
		},
		sortType: {
			type: String,
			require: true,
			default: 'desc'
		},
		pageSize: {
			type: Number,
			require: true,
			default: 3
		}
	},
	data() {
		return {
			pageNumber:1,
			newsList: [],
			displayedList:[],
			tagColors: {
				typeScript: '#3178c6',
				javaScript: '#f1e05a',
				html: '#e34c26',
				css: '#563d7c',
				java: 'orange',
				python: '#3572A5',
				golang: '#00ADD8',
				go: '#00ADD8',
				shell: '#89e051',
				'c++': '#f34b7d',
				other: '#ededed'
			},
			loading: '',
		}
	},
	created() {
		this.getnews();
	},
	methods: {
		handleNumber(number) {
			if (number >= 1000) {
				number = number / 1000;
				return Math.round(number) + "k";
			}
			return number;
		},
		handleTagColor(language) {
			let color = this.tagColors['other'];
			if (!language) {
				return color;
			}
			Object.keys(this.tagColors).forEach(key => {
				if (key.toLowerCase() === language.toLowerCase()) {
					color = this.tagColors[key];
				}
			});
			return color;
		},
		async getnews() {
			const url = "https://newsapi.org/v2/top-headlines?country=us&apiKey=804f6bdac084463ba3fadb53f9efce90";
			let response = await this.$axios.get(url);
			this.newsList = response.articles;
			this.displayedList = this.newsList.slice(0,5)
			console.log(this.newsList);
		},
		async getMorenews() {
			this.loading = 'el-icon-loading';
			const l = Math.min(this.displayedList.length + 3, this.newsList.length);
			this.displayedList = this.newsList.slice(0, l);
			this.loading = '';
		},
	}
}
</script>

<style lang="less" scoped>
.box-card {
	width: 100%;

	.title {
		background-image: linear-gradient(0deg,
				rgba(0, 0, 0, 0.3) 0,
				transparent);
		background-color: #0079d3;
		height: 80px;
		width: 100%;
		color: #fff;
		font-size: 20px;
		line-height: 80px;
		padding-left: 10px;
		box-sizing: border-box;
		text-align: center;
		border-radius: 4px 4px 0px 0px;
	}

	.trending-news-list {
		.trending-news-item {
			margin-top: 1rem;
			padding: 0px 5px 0px 5px;

			.user-info {
				display: flex;

				.avatar {
					margin-right: 10px;
				}

				.news-name {
					font-weight: 600;
					display: flex;
					align-items: center;

					a {
						white-space: nowrap;
						overflow: hidden;
						text-overflow: ellipsis;
					}
				}
			}

			.introduction {
				margin: 0.5rem 0;
				font-size: 12px;
				display: flex;
				gap: 8px;
				align-items: center;
			}

			.rank {
				font-weight: bold;
				color: #555;
			}

			.meta {
				font-size: 12px;
				color: #a7a3a3;

				i {
					margin-right: 2px;
				}

				.forks {
					margin: 0 8px;
				}
			}
		}

		.next-btn {
			background-image: linear-gradient(0deg,
					rgba(0, 0, 0, 0.3) 0,
					transparent);
			background-color: #0079d3;
			color: #fff;
			height: 40px;
			line-height: 40px;
			text-align: center;
			cursor: pointer;
			font-weight: 600;

			i {
				margin-right: 5px;
			}
		}
	}
}
</style>
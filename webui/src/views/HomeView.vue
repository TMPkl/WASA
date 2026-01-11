<script>
import { logout } from '../services/auth.js'
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		},
		logout() {
			logout()
			this.$router.replace('/login')
		},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">My Messages</h1>
			<div class="btn-toolbar toolbar-button"> 
				<div class="btn-group">
					<button type="button" class="refresh-btn btn-bg  me-2" @click="refresh" :disabled="loading">
						<span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
						Refresh
					</button>
					<button type="button" class="logout-btn btn-bg  me-2" :disabled="loading" @click="logout">
						Logout
					</button>
				
					<button type="button" class=" nc-btn btn-bg " @click="newItem">
						New Conversation
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
	.toolbar-button {
		margin-right: 8px;
	}
	.logout-btn {background: #ff1509;
				 color:white ;
				 border-color: #ff1509;
				 border-width: 1px;
				 border-style: solid;
				}
		.logout-btn:hover {background: white;
					  		color:#ff1509 ;}
	.refresh-btn {background:#0dc01c ;
				 color: white;
				 border-color: #0dc01c;
				 border-width: 1px;
				 border-style: solid;
				border-top-left-radius: 6px;
				border-bottom-left-radius: 6px;
				}
		.refresh-btn:hover {background: white;
					  		color: #0dc01c;}
	.nc-btn {background: #007bff;
				color: white;
				border-color: #007bff;
				border-width: 1px;
				border-style: solid;
				border-top-right-radius: 6px;
				border-bottom-right-radius: 6px;
			}
		.nc-btn:hover {background: white;
						color: #0056b3;}	

	.refresh-btn, .logout-btn, .nc-btn {
		font-size: 1.15rem;
		padding: 0.65rem 1rem;
		line-height: 1.4;
		border-radius: 0.6rem;
	}

</style>

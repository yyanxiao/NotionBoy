(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[459],{4699:function(e,t,s){(window.__NEXT_P=window.__NEXT_P||[]).push(["/login",function(){return s(5887)}])},5887:function(e,t,s){"use strict";s.r(t),s.d(t,{default:function(){return p}});var n=s(1527),r=s(3933),i=s(4271),l=s(959),a=s(9140);let o=l.forwardRef((e,t)=>{let{className:s,...r}=e;return(0,n.jsx)("input",{className:(0,a.cn)("flex h-10 w-full rounded-md border border-slate-300 bg-transparent py-2 px-3 text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:text-slate-50 dark:focus:ring-slate-400 dark:focus:ring-offset-slate-900",s),ref:t,...r})});o.displayName="Input";var c=s(9300),d=s(4379),u=s(4157),f=s(6131),h=s(7765),x=s.n(h),m=s(3125);function p(){let e=(0,m.useRouter)(),{toast:t}=(0,d.pm)(),[s,a]=(0,l.useState)([]),[h,p]=(0,l.useState)(!1),[g,w]=(0,l.useState)(""),N=[{id:"github",name:"Github"}];(0,l.useEffect)(()=>{v(),h||N.forEach(async e=>{let t=await u.t.OAuthURL({provider:e.id});e.url=t.url,a([...s,e])})},[]);let v=()=>{let e=f.Z.get(r.o3);if(void 0==e){p(!1);return}p(!0)};if(h)return(0,n.jsx)("div",{className:"flex-grow container mx-auto flex flex-col p-8 items-center",children:(0,n.jsx)("div",{className:"prose m-4",children:(0,n.jsx)("h1",{children:"You are already logged in!"})})});let j=s=>{u.t.GenrateToken({magicCode:s}).then(t=>{f.Z.set(r.o3,t.token),f.Z.set(r.kk,t.expiry),p(!0),e.push(c.J.links.home)}).catch(e=>{t({variant:"destructive",title:"Login with MagicCode error",description:JSON.stringify(e)})})};return(0,n.jsxs)("div",{className:"flex-grow container mx-auto flex flex-col p-8 items-center",children:[(0,n.jsx)("div",{className:"prose m-4",children:(0,n.jsx)("h1",{children:"Please Login"})}),(0,n.jsxs)("div",{className:"container mx-auto max-w-sm flex flex-col bg-white p-4 items-center",children:[s.map(e=>void 0===e.url?(0,n.jsxs)(i.z,{disabled:!0,children:["Sign In with ",e.name]},e.id):(0,n.jsx)(x(),{href:e.url,className:"w-full m-2",children:(0,n.jsxs)(i.z,{className:"w-full",children:["Sign In with ",e.name]})},e.id)),(0,n.jsxs)("div",{className:"flex w-full max-w-sm items-center space-x-2",children:[(0,n.jsx)(o,{type:"text",placeholder:"Login with MagicCode",value:g,onChange:e=>w(e.target.value)}),(0,n.jsx)(i.z,{type:"submit",onClick:()=>j(g),children:"Login"})]})]})]})}}},function(e){e.O(0,[774,888,179],function(){return e(e.s=4699)}),_N_E=e.O()}]);
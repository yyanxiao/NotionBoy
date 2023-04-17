"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[933],{3933:function(t,e,r){r.d(e,{c_:function(){return f},kk:function(){return h},o3:function(){return l}});var o=r(11527),n=r(49300),a=r(4379),s=r(44157),i=r(96131),c=r(73125),d=r(50959),u=r(64271);let l="token",h="tokenExpire",p=n.J.authPages;function f(){let[t,e]=(0,d.useState)(null),[r,f]=(0,d.useState)(""),[g,v]=(0,d.useState)(!0),m=(0,c.useRouter)(),{toast:b}=(0,a.pm)();(0,d.useEffect)(()=>{if(m.isReady){m.pathname==n.J.links.login&&v(!1);let t=i.Z.get(l);t&&e(t),y()}},[m]),(0,d.useEffect)(()=>{r&&b({variant:"destructive",title:"Auth error",description:r})},[r]);let y=()=>{for(let t of p)if(m.pathname.startsWith(t)){let t=i.Z.get(h),r=Date.parse(t)-Date.now();r<0?(i.Z.remove(l),i.Z.remove(h)):r<3600?s.t.GenrateToken({}).then(t=>{t.token?(i.Z.set(l,t.token,{expires:new Date(t.expiry)}),i.Z.set(h,t.expiry,{expires:new Date(t.expiry)}),e(t.token)):(f("Get token failed"),S())}).catch(t=>{f("Get token failed ".concat(t)),S()}):e(i.Z.get(l))}},k=()=>{i.Z.remove(l),i.Z.remove(h),m.push(n.J.links.home),m.reload()},S=()=>{m.push(n.J.links.login)};return(0,o.jsx)(o.Fragment,{children:(()=>{if(g)return t?(0,o.jsx)(u.z,{variant:"ghost",size:"sm",onClick:k,children:"Logout"}):(0,o.jsx)(u.z,{variant:"ghost",size:"sm",onClick:S,children:"Login"})})()})}},64271:function(t,e,r){r.d(e,{d:function(){return i},z:function(){return c}});var o=r(11527),n=r(50959),a=r(71590),s=r(29140);let i=(0,a.j)("active:scale-95 inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 dark:hover:bg-slate-800 dark:hover:text-slate-100 disabled:opacity-50 dark:focus:ring-slate-400 disabled:pointer-events-none dark:focus:ring-offset-slate-900 data-[state=open]:bg-slate-100 dark:data-[state=open]:bg-slate-800",{variants:{variant:{default:"bg-slate-900 text-white hover:bg-slate-700 dark:bg-slate-50 dark:text-slate-900",destructive:"bg-red-500 text-white hover:bg-red-600 dark:hover:bg-red-600",outline:"bg-transparent border border-slate-200 hover:bg-slate-100 dark:border-slate-700 dark:text-slate-100",subtle:"bg-slate-100 text-slate-900 hover:bg-slate-200 dark:bg-slate-700 dark:text-slate-100",ghost:"bg-transparent hover:bg-slate-100 dark:hover:bg-slate-800 dark:text-slate-100 dark:hover:text-slate-100 data-[state=open]:bg-transparent dark:data-[state=open]:bg-transparent",link:"bg-transparent dark:bg-transparent underline-offset-4 hover:underline text-slate-900 dark:text-slate-100 hover:bg-transparent dark:hover:bg-transparent"},size:{default:"h-10 py-2 px-4",sm:"h-9 px-2 rounded-md",lg:"h-11 px-8 rounded-md"}},defaultVariants:{variant:"default",size:"default"}}),c=n.forwardRef((t,e)=>{let{className:r,variant:n,size:a,...c}=t;return(0,o.jsx)("button",{className:(0,s.cn)(i({variant:n,size:a,className:r})),ref:e,...c})});c.displayName="Button"},49300:function(t,e,r){function o(t){return t+".html"}r.d(e,{J:function(){return n}});let n={name:"NotionBoy",description:"NotionBoy is a note app base on Notion. It's a web app, you can use it in your browser.",links:{twitter:"https://twitter.com/LiuVaayne",github:"https://github.com/vaayne/NotionBoy",chatgpt:o("/chat"),login:o("/login"),home:"/",authCallback:o("/authcallback"),price:o("/price"),order:o("/order"),profile:o("/profile")},authPages:[o("/chat"),o("/order"),o("/user")]}},44157:function(t,e,r){r.d(e,{t:function(){return u}});let o=Array(64),n=Array(123);for(let t=0;t<64;)n[o[t]=t<26?t+65:t<52?t+71:t<62?t-4:t-59|43]=t++;function a(t,e){return e&&e.constructor===Uint8Array?function(t,e,r){let n=null,a=[],s=0,i=0,c;for(;e<r;){let r=t[e++];switch(i){case 0:a[s++]=o[r>>2],c=(3&r)<<4,i=1;break;case 1:a[s++]=o[c|r>>4],c=(15&r)<<2,i=2;break;case 2:a[s++]=o[c|r>>6],a[s++]=o[63&r],i=0}s>8191&&((n||(n=[])).push(String.fromCharCode.apply(String,a)),s=0)}return(i&&(a[s++]=o[c],a[s++]=61,1===i&&(a[s++]=61)),n)?(s&&n.push(String.fromCharCode.apply(String,a.slice(0,s))),n.join("")):String.fromCharCode.apply(String,a.slice(0,s))}(e,0,e.length):e}function s(t,e){let{pathPrefix:r,...o}=e||{},n=r?"".concat(r).concat(t):t;return fetch(n,o).then(t=>t.json().then(e=>{if(!t.ok)throw e;return e}))}async function i(t,e,r){var o;let{pathPrefix:n,...a}=r||{},s=n?"".concat(n).concat(t):t,i=await fetch(s,a);if(!i.ok){let t=await i.json(),e=t.error&&t.error.message?t.error.message:"";throw Error(e)}if(!i.body)throw Error("response doesnt have a body");await i.body.pipeThrough(new TextDecoderStream).pipeThrough(new TransformStream({start(t){t.buf="",t.pos=0},transform(t,e){for(void 0===e.buf&&(e.buf=""),void 0===e.pos&&(e.pos=0),e.buf+=t;e.pos<e.buf.length;)if("\n"===e.buf[e.pos]){let t=e.buf.substring(0,e.pos),r=JSON.parse(t);e.enqueue(r.result),e.buf=e.buf.substring(e.pos+1),e.pos=0}else++e.pos}})).pipeTo((o=t=>{e&&e(t)},new WritableStream({write(t){o(t)}})))}function c(t){return["string","number","boolean"].some(e=>typeof t===e)}function d(t){let e=arguments.length>1&&void 0!==arguments[1]?arguments[1]:[],r=function t(e){let r=arguments.length>1&&void 0!==arguments[1]?arguments[1]:"";return Object.keys(e).reduce((o,n)=>{let a=e[n],s=r?[r,n].join("."):n,i=Array.isArray(a)&&a.every(t=>c(t))&&a.length>0,d=c(a)&&!(!1===a||0===a||""===a),u={};return!function(t){let e="Object"===Object.prototype.toString.call(t).slice(8,-1);if(!(null!==t&&e)||!e)return!1;let r=Object.getPrototypeOf(t),o="object"==typeof r&&r.constructor===Object.prototype.constructor;return o}(a)?(d||i)&&(u={[s]:a}):u=t(a,s),{...o,...u}},{})}(t),o=Object.keys(r).reduce((t,o)=>{let n=r[o];return e.find(t=>t===o)?t:Array.isArray(n)?[...t,...n.map(t=>[o,t.toString()])]:t=[...t,[o,n.toString()]]},[]);return new URLSearchParams(o).toString()}class u{static Status(t,e){return s("/v1/status?".concat(d(t,[])),{...e,method:"GET"})}static GenrateToken(t,e){return s("/v1/auth/token",{...e,method:"POST",body:JSON.stringify(t,a)})}static OAuthProviders(t,e){return s("/v1/auth/providers?".concat(d(t,[])),{...e,method:"GET"})}static OAuthCallback(t,e){return s("/v1/auth/callback",{...e,method:"POST",body:JSON.stringify(t,a)})}static GenerateApiKey(t,e){return s("/v1/auth/apikey",{...e,method:"POST"})}static DeleteApiKey(t,e){return s("/v1/auth/apikey",{...e,method:"DELETE"})}static GenerateWechatQRCode(t,e){return s("/v1/auth/wechat/qrcode",{...e,method:"POST",body:JSON.stringify(t,a)})}static CreateConversation(t,e){return s("/v1/conversations",{...e,method:"POST",body:JSON.stringify(t,a)})}static UpdateConversation(t,e){return s("/v1/conversations/".concat(t.id),{...e,method:"PUT",body:JSON.stringify(t,a)})}static GetConversation(t,e){return s("/v1/conversations/".concat(t.id,"?").concat(d(t,["id"])),{...e,method:"GET"})}static ListConversations(t,e){return s("/v1/conversations?".concat(d(t,[])),{...e,method:"GET"})}static DeleteConversation(t,e){return s("/v1/conversations/".concat(t.id),{...e,method:"DELETE"})}static CreateMessage(t,e,r){return i("/v1/conversations/".concat(t.conversationId,"/messages"),e,{...r,method:"POST",body:JSON.stringify(t,a)})}static UpdateMessage(t,e,r){return i("/v1/conversations/".concat(t.conversationId,"/messages/").concat(t.id),e,{...r,method:"POST",body:JSON.stringify(t,a)})}static GetMessage(t,e){return s("/v1/conversations/".concat(t.conversationId,"/messages/").concat(t.id,"?").concat(d(t,["conversationId","id"])),{...e,method:"GET"})}static ListMessages(t,e){return s("/v1/conversations/".concat(t.conversationId,"/messages?").concat(d(t,["conversationId"])),{...e,method:"GET"})}static DeleteMessage(t,e){return s("/v1/conversations/".concat(t.conversationId,"/messages/").concat(t.id),{...e,method:"DELETE"})}static CreateOrder(t,e){return s("/v1/orders",{...e,method:"POST",body:JSON.stringify(t,a)})}static GetOrder(t,e){return s("/v1/orders/".concat(t.id,"?").concat(d(t,["id"])),{...e,method:"GET"})}static ListOrders(t,e){return s("/v1/orders?".concat(d(t,[])),{...e,method:"GET"})}static DeleteOrder(t,e){return s("/v1/orders/".concat(t.id),{...e,method:"DELETE"})}static UpdateOrder(t,e){return s("/v1/orders/".concat(t.id),{...e,method:"PATCH",body:JSON.stringify(t,a)})}static PayOrder(t,e){return s("/v1/orders/".concat(t.id,"/pay"),{...e,method:"POST",body:JSON.stringify(t,a)})}static CreateProduct(t,e){return s("/v1/products",{...e,method:"POST",body:JSON.stringify(t,a)})}static GetProduct(t,e){return s("/v1/products/".concat(t.id,"?").concat(d(t,["id"])),{...e,method:"GET"})}static ListProducts(t,e){return s("/v1/products?".concat(d(t,[])),{...e,method:"GET"})}static DeleteProduct(t,e){return s("/v1/products/".concat(t.id),{...e,method:"DELETE"})}static UpdateProduct(t,e){return s("/v1/products/".concat(t.id),{...e,method:"PATCH",body:JSON.stringify(t,a)})}static ListPrompts(t,e){return s("/v1/prompts?".concat(d(t,[])),{...e,method:"GET"})}static GetPrompt(t,e){return s("/v1/prompts/".concat(t.id,"?").concat(d(t,["id"])),{...e,method:"GET"})}static CreatePrompt(t,e){return s("/v1/prompts",{...e,method:"POST",body:JSON.stringify(t,a)})}static UpdatePrompt(t,e){return s("/v1/prompts/".concat(t.id),{...e,method:"PATCH",body:JSON.stringify(t,a)})}static DeletePrompt(t,e){return s("/v1/prompts/".concat(t.id),{...e,method:"DELETE"})}}}}]);
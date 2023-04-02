"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[84],{2084:function(e,t,n){n.d(t,{Dx:function(){return G},VY:function(){return q},aV:function(){return U},dk:function(){return J},fC:function(){return Y},h_:function(){return H},x8:function(){return Q},xz:function(){return z}});var r=n(1163),o=n(959),a=n(6459),c=n(7381),u=n(780),i=n(7055),l=n(5087),d=n(5886),s=n(5043),f=n(2339),p=n(9202),v=n(8975),m=n(1912),h=n(434),g=n(9343),E=n(7e3);let y="Dialog",[b,w]=(0,u.b)(y),[C,R]=b(y),D=e=>{let{__scopeDialog:t,children:n,open:r,defaultOpen:a,onOpenChange:c,modal:u=!0}=e,d=(0,o.useRef)(null),s=(0,o.useRef)(null),[f=!1,p]=(0,l.T)({prop:r,defaultProp:a,onChange:c});return(0,o.createElement)(C,{scope:t,triggerRef:d,contentRef:s,contentId:(0,i.M)(),titleId:(0,i.M)(),descriptionId:(0,i.M)(),open:f,onOpenChange:p,onOpenToggle:(0,o.useCallback)(()=>p(e=>!e),[p]),modal:u},n)},M=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,...u}=e,i=R("DialogTrigger",n),l=(0,c.e)(t,i.triggerRef);return(0,o.createElement)(v.WV.button,(0,r.Z)({type:"button","aria-haspopup":"dialog","aria-expanded":i.open,"aria-controls":i.contentId,"data-state":K(i.open)},u,{ref:l,onClick:(0,a.M)(e.onClick,i.onOpenToggle)}))}),S="DialogPortal",[k,A]=b(S,{forceMount:void 0}),_=e=>{let{__scopeDialog:t,forceMount:n,children:r,container:a}=e,c=R(S,t);return(0,o.createElement)(k,{scope:t,forceMount:n},o.Children.map(r,e=>(0,o.createElement)(p.z,{present:n||c.open},(0,o.createElement)(f.h,{asChild:!0,container:a},e))))},T="DialogOverlay",O=(0,o.forwardRef)((e,t)=>{let n=A(T,e.__scopeDialog),{forceMount:a=n.forceMount,...c}=e,u=R(T,e.__scopeDialog);return u.modal?(0,o.createElement)(p.z,{present:a||u.open},(0,o.createElement)(x,(0,r.Z)({},c,{ref:t}))):null}),x=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,...a}=e,c=R(T,n);return(0,o.createElement)(h.Z,{as:E.g7,allowPinchZoom:!0,shards:[c.contentRef]},(0,o.createElement)(v.WV.div,(0,r.Z)({"data-state":K(c.open)},a,{ref:t,style:{pointerEvents:"auto",...a.style}})))}),N="DialogContent",I=(0,o.forwardRef)((e,t)=>{let n=A(N,e.__scopeDialog),{forceMount:a=n.forceMount,...c}=e,u=R(N,e.__scopeDialog);return(0,o.createElement)(p.z,{present:a||u.open},u.modal?(0,o.createElement)(L,(0,r.Z)({},c,{ref:t})):(0,o.createElement)(P,(0,r.Z)({},c,{ref:t})))}),L=(0,o.forwardRef)((e,t)=>{let n=R(N,e.__scopeDialog),u=(0,o.useRef)(null),i=(0,c.e)(t,n.contentRef,u);return(0,o.useEffect)(()=>{let e=u.current;if(e)return(0,g.Ry)(e)},[]),(0,o.createElement)(F,(0,r.Z)({},e,{ref:i,trapFocus:n.open,disableOutsidePointerEvents:!0,onCloseAutoFocus:(0,a.M)(e.onCloseAutoFocus,e=>{var t;e.preventDefault(),null===(t=n.triggerRef.current)||void 0===t||t.focus()}),onPointerDownOutside:(0,a.M)(e.onPointerDownOutside,e=>{let t=e.detail.originalEvent,n=0===t.button&&!0===t.ctrlKey,r=2===t.button||n;r&&e.preventDefault()}),onFocusOutside:(0,a.M)(e.onFocusOutside,e=>e.preventDefault())}))}),P=(0,o.forwardRef)((e,t)=>{let n=R(N,e.__scopeDialog),a=(0,o.useRef)(!1);return(0,o.createElement)(F,(0,r.Z)({},e,{ref:t,trapFocus:!1,disableOutsidePointerEvents:!1,onCloseAutoFocus:t=>{var r,o;null===(r=e.onCloseAutoFocus)||void 0===r||r.call(e,t),t.defaultPrevented||(a.current||null===(o=n.triggerRef.current)||void 0===o||o.focus(),t.preventDefault()),a.current=!1},onInteractOutside:t=>{var r,o;null===(r=e.onInteractOutside)||void 0===r||r.call(e,t),t.defaultPrevented||(a.current=!0);let c=t.target,u=null===(o=n.triggerRef.current)||void 0===o?void 0:o.contains(c);u&&t.preventDefault()}}))}),F=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,trapFocus:a,onOpenAutoFocus:u,onCloseAutoFocus:i,...l}=e,f=R(N,n),p=(0,o.useRef)(null),v=(0,c.e)(t,p);return(0,m.EW)(),(0,o.createElement)(o.Fragment,null,(0,o.createElement)(s.M,{asChild:!0,loop:!0,trapped:a,onMountAutoFocus:u,onUnmountAutoFocus:i},(0,o.createElement)(d.XB,(0,r.Z)({role:"dialog",id:f.contentId,"aria-describedby":f.descriptionId,"aria-labelledby":f.titleId,"data-state":K(f.open)},l,{ref:v,onDismiss:()=>f.onOpenChange(!1)}))),!1)}),W="DialogTitle",Z=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,...a}=e,c=R(W,n);return(0,o.createElement)(v.WV.h2,(0,r.Z)({id:c.titleId},a,{ref:t}))}),B=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,...a}=e,c=R("DialogDescription",n);return(0,o.createElement)(v.WV.p,(0,r.Z)({id:c.descriptionId},a,{ref:t}))}),j=(0,o.forwardRef)((e,t)=>{let{__scopeDialog:n,...c}=e,u=R("DialogClose",n);return(0,o.createElement)(v.WV.button,(0,r.Z)({type:"button"},c,{ref:t,onClick:(0,a.M)(e.onClick,()=>u.onOpenChange(!1))}))});function K(e){return e?"open":"closed"}let[V,X]=(0,u.k)("DialogTitleWarning",{contentName:N,titleName:W,docsSlug:"dialog"}),Y=D,z=M,H=_,U=O,q=I,G=Z,J=B,Q=j},1912:function(e,t,n){n.d(t,{EW:function(){return a}});var r=n(959);let o=0;function a(){(0,r.useEffect)(()=>{var e,t;let n=document.querySelectorAll("[data-radix-focus-guard]");return document.body.insertAdjacentElement("afterbegin",null!==(e=n[0])&&void 0!==e?e:c()),document.body.insertAdjacentElement("beforeend",null!==(t=n[1])&&void 0!==t?t:c()),o++,()=>{1===o&&document.querySelectorAll("[data-radix-focus-guard]").forEach(e=>e.remove()),o--}},[])}function c(){let e=document.createElement("span");return e.setAttribute("data-radix-focus-guard",""),e.tabIndex=0,e.style.cssText="outline: none; opacity: 0; position: fixed; pointer-events: none",e}},5043:function(e,t,n){let r;n.d(t,{M:function(){return f}});var o=n(1163),a=n(959),c=n(7381),u=n(8975),i=n(9928);let l="focusScope.autoFocusOnMount",d="focusScope.autoFocusOnUnmount",s={bubbles:!1,cancelable:!0},f=(0,a.forwardRef)((e,t)=>{let{loop:n=!1,trapped:r=!1,onMountAutoFocus:f,onUnmountAutoFocus:g,...E}=e,[y,b]=(0,a.useState)(null),w=(0,i.W)(f),C=(0,i.W)(g),R=(0,a.useRef)(null),D=(0,c.e)(t,e=>b(e)),M=(0,a.useRef)({paused:!1,pause(){this.paused=!0},resume(){this.paused=!1}}).current;(0,a.useEffect)(()=>{if(r){function e(e){if(M.paused||!y)return;let t=e.target;y.contains(t)?R.current=t:m(R.current,{select:!0})}function t(e){M.paused||!y||y.contains(e.relatedTarget)||m(R.current,{select:!0})}return document.addEventListener("focusin",e),document.addEventListener("focusout",t),()=>{document.removeEventListener("focusin",e),document.removeEventListener("focusout",t)}}},[r,y,M.paused]),(0,a.useEffect)(()=>{if(y){h.add(M);let e=document.activeElement,t=y.contains(e);if(!t){let t=new CustomEvent(l,s);y.addEventListener(l,w),y.dispatchEvent(t),t.defaultPrevented||(function(e,{select:t=!1}={}){let n=document.activeElement;for(let r of e)if(m(r,{select:t}),document.activeElement!==n)return}(p(y).filter(e=>"A"!==e.tagName),{select:!0}),document.activeElement===e&&m(y))}return()=>{y.removeEventListener(l,w),setTimeout(()=>{let t=new CustomEvent(d,s);y.addEventListener(d,C),y.dispatchEvent(t),t.defaultPrevented||m(null!=e?e:document.body,{select:!0}),y.removeEventListener(d,C),h.remove(M)},0)}}},[y,w,C,M]);let S=(0,a.useCallback)(e=>{if(!n&&!r||M.paused)return;let t="Tab"===e.key&&!e.altKey&&!e.ctrlKey&&!e.metaKey,o=document.activeElement;if(t&&o){let t=e.currentTarget,[r,a]=function(e){let t=p(e),n=v(t,e),r=v(t.reverse(),e);return[n,r]}(t);r&&a?e.shiftKey||o!==a?e.shiftKey&&o===r&&(e.preventDefault(),n&&m(a,{select:!0})):(e.preventDefault(),n&&m(r,{select:!0})):o===t&&e.preventDefault()}},[n,r,M.paused]);return(0,a.createElement)(u.WV.div,(0,o.Z)({tabIndex:-1},E,{ref:D,onKeyDown:S}))});function p(e){let t=[],n=document.createTreeWalker(e,NodeFilter.SHOW_ELEMENT,{acceptNode:e=>{let t="INPUT"===e.tagName&&"hidden"===e.type;return e.disabled||e.hidden||t?NodeFilter.FILTER_SKIP:e.tabIndex>=0?NodeFilter.FILTER_ACCEPT:NodeFilter.FILTER_SKIP}});for(;n.nextNode();)t.push(n.currentNode);return t}function v(e,t){for(let n of e)if(!function(e,{upTo:t}){if("hidden"===getComputedStyle(e).visibility)return!0;for(;e&&(void 0===t||e!==t);){if("none"===getComputedStyle(e).display)return!0;e=e.parentElement}return!1}(n,{upTo:t}))return n}function m(e,{select:t=!1}={}){if(e&&e.focus){var n;let r=document.activeElement;e.focus({preventScroll:!0}),e!==r&&(n=e)instanceof HTMLInputElement&&"select"in n&&t&&e.select()}}let h=(r=[],{add(e){let t=r[0];e!==t&&(null==t||t.pause()),(r=g(r,e)).unshift(e)},remove(e){var t;null===(t=(r=g(r,e))[0])||void 0===t||t.resume()}});function g(e,t){let n=[...e],r=n.indexOf(t);return -1!==r&&n.splice(r,1),n}},9343:function(e,t,n){n.d(t,{Ry:function(){return l}});var r=new WeakMap,o=new WeakMap,a={},c=0,u=function(e){return e&&(e.host||u(e.parentNode))},i=function(e,t,n,i){var l=(Array.isArray(e)?e:[e]).map(function(e){if(t.contains(e))return e;var n=u(e);return n&&t.contains(n)?n:(console.error("aria-hidden",e,"in not contained inside",t,". Doing nothing"),null)}).filter(function(e){return Boolean(e)});a[n]||(a[n]=new WeakMap);var d=a[n],s=[],f=new Set,p=new Set(l),v=function(e){!e||f.has(e)||(f.add(e),v(e.parentNode))};l.forEach(v);var m=function(e){!e||p.has(e)||Array.prototype.forEach.call(e.children,function(e){if(f.has(e))m(e);else{var t=e.getAttribute(i),a=null!==t&&"false"!==t,c=(r.get(e)||0)+1,u=(d.get(e)||0)+1;r.set(e,c),d.set(e,u),s.push(e),1===c&&a&&o.set(e,!0),1===u&&e.setAttribute(n,"true"),a||e.setAttribute(i,"true")}})};return m(t),f.clear(),c++,function(){s.forEach(function(e){var t=r.get(e)-1,a=d.get(e)-1;r.set(e,t),d.set(e,a),t||(o.has(e)||e.removeAttribute(i),o.delete(e)),a||e.removeAttribute(n)}),--c||(r=new WeakMap,r=new WeakMap,o=new WeakMap,a={})}},l=function(e,t,n){void 0===n&&(n="data-aria-hidden");var r=Array.from(Array.isArray(e)?e:[e]),o=t||("undefined"==typeof document?null:(Array.isArray(e)?e[0]:e).ownerDocument.body);return o?(r.push.apply(r,Array.from(o.querySelectorAll("[aria-live]"))),i(r,o,n,"aria-hidden")):function(){return null}}},434:function(e,t,n){n.d(t,{Z:function(){return V}});var r,o,a,c,u,i,l=function(){return(l=Object.assign||function(e){for(var t,n=1,r=arguments.length;n<r;n++)for(var o in t=arguments[n])Object.prototype.hasOwnProperty.call(t,o)&&(e[o]=t[o]);return e}).apply(this,arguments)};function d(e,t){var n={};for(var r in e)Object.prototype.hasOwnProperty.call(e,r)&&0>t.indexOf(r)&&(n[r]=e[r]);if(null!=e&&"function"==typeof Object.getOwnPropertySymbols)for(var o=0,r=Object.getOwnPropertySymbols(e);o<r.length;o++)0>t.indexOf(r[o])&&Object.prototype.propertyIsEnumerable.call(e,r[o])&&(n[r[o]]=e[r[o]]);return n}var s=n(959),f="right-scroll-bar-position",p="width-before-scroll-bar",v=(void 0===r&&(r={}),(void 0===o&&(o=function(e){return e}),a=[],c=!1,u={read:function(){if(c)throw Error("Sidecar: could not `read` from an `assigned` medium. `read` could be used only with `useMedium`.");return a.length?a[a.length-1]:null},useMedium:function(e){var t=o(e,c);return a.push(t),function(){a=a.filter(function(e){return e!==t})}},assignSyncMedium:function(e){for(c=!0;a.length;){var t=a;a=[],t.forEach(e)}a={push:function(t){return e(t)},filter:function(){return a}}},assignMedium:function(e){c=!0;var t=[];if(a.length){var n=a;a=[],n.forEach(e),t=a}var r=function(){var n=t;t=[],n.forEach(e)},o=function(){return Promise.resolve().then(r)};o(),a={push:function(e){t.push(e),o()},filter:function(e){return t=t.filter(e),a}}}}).options=l({async:!0,ssr:!1},r),u),m=function(){},h=s.forwardRef(function(e,t){var n,r,o,a=s.useRef(null),c=s.useState({onScrollCapture:m,onWheelCapture:m,onTouchMoveCapture:m}),u=c[0],i=c[1],f=e.forwardProps,p=e.children,h=e.className,g=e.removeScrollBar,E=e.enabled,y=e.shards,b=e.sideCar,w=e.noIsolation,C=e.inert,R=e.allowPinchZoom,D=e.as,M=d(e,["forwardProps","children","className","removeScrollBar","enabled","shards","sideCar","noIsolation","inert","allowPinchZoom","as"]),S=(n=[a,t],r=function(e){return n.forEach(function(t){var n;return"function"==typeof(n=t)?n(e):n&&(n.current=e),n})},(o=(0,s.useState)(function(){return{value:null,callback:r,facade:{get current(){return o.value},set current(value){var e=o.value;e!==value&&(o.value=value,o.callback(value,e))}}}})[0]).callback=r,o.facade),k=l(l({},M),u);return s.createElement(s.Fragment,null,E&&s.createElement(b,{sideCar:v,removeScrollBar:g,shards:y,noIsolation:w,inert:C,setCallbacks:i,allowPinchZoom:!!R,lockRef:a}),f?s.cloneElement(s.Children.only(p),l(l({},k),{ref:S})):s.createElement(void 0===D?"div":D,l({},k,{className:h,ref:S}),p))});h.defaultProps={enabled:!0,removeScrollBar:!0,inert:!1},h.classNames={fullWidth:p,zeroRight:f};var g=function(e){var t=e.sideCar,n=d(e,["sideCar"]);if(!t)throw Error("Sidecar: please provide `sideCar` property to import the right car");var r=t.read();if(!r)throw Error("Sidecar medium not found");return s.createElement(r,l({},n))};g.isSideCarExport=!0;var E=function(){var e=0,t=null;return{add:function(r){if(0==e&&(t=function(){if(!document)return null;var e=document.createElement("style");e.type="text/css";var t=i||n.nc;return t&&e.setAttribute("nonce",t),e}())){var o,a;(o=t).styleSheet?o.styleSheet.cssText=r:o.appendChild(document.createTextNode(r)),a=t,(document.head||document.getElementsByTagName("head")[0]).appendChild(a)}e++},remove:function(){--e||!t||(t.parentNode&&t.parentNode.removeChild(t),t=null)}}},y=function(){var e=E();return function(t,n){s.useEffect(function(){return e.add(t),function(){e.remove()}},[t&&n])}},b=function(){var e=y();return function(t){return e(t.styles,t.dynamic),null}},w={left:0,top:0,right:0,gap:0},C=function(e){return parseInt(e||"",10)||0},R=function(e){var t=window.getComputedStyle(document.body),n=t["padding"===e?"paddingLeft":"marginLeft"],r=t["padding"===e?"paddingTop":"marginTop"],o=t["padding"===e?"paddingRight":"marginRight"];return[C(n),C(r),C(o)]},D=function(e){if(void 0===e&&(e="margin"),"undefined"==typeof window)return w;var t=R(e),n=document.documentElement.clientWidth,r=window.innerWidth;return{left:t[0],top:t[1],right:t[2],gap:Math.max(0,r-n+t[2]-t[0])}},M=b(),S=function(e,t,n,r){var o=e.left,a=e.top,c=e.right,u=e.gap;return void 0===n&&(n="margin"),"\n  .".concat("with-scroll-bars-hidden"," {\n   overflow: hidden ").concat(r,";\n   padding-right: ").concat(u,"px ").concat(r,";\n  }\n  body {\n    overflow: hidden ").concat(r,";\n    overscroll-behavior: contain;\n    ").concat([t&&"position: relative ".concat(r,";"),"margin"===n&&"\n    padding-left: ".concat(o,"px;\n    padding-top: ").concat(a,"px;\n    padding-right: ").concat(c,"px;\n    margin-left:0;\n    margin-top:0;\n    margin-right: ").concat(u,"px ").concat(r,";\n    "),"padding"===n&&"padding-right: ".concat(u,"px ").concat(r,";")].filter(Boolean).join(""),"\n  }\n  \n  .").concat(f," {\n    right: ").concat(u,"px ").concat(r,";\n  }\n  \n  .").concat(p," {\n    margin-right: ").concat(u,"px ").concat(r,";\n  }\n  \n  .").concat(f," .").concat(f," {\n    right: 0 ").concat(r,";\n  }\n  \n  .").concat(p," .").concat(p," {\n    margin-right: 0 ").concat(r,";\n  }\n  \n  body {\n    ").concat("--removed-body-scroll-bar-size",": ").concat(u,"px;\n  }\n")},k=function(e){var t=e.noRelative,n=e.noImportant,r=e.gapMode,o=void 0===r?"margin":r,a=s.useMemo(function(){return D(o)},[o]);return s.createElement(M,{styles:S(a,!t,o,n?"":"!important")})},A=!1;if("undefined"!=typeof window)try{var _=Object.defineProperty({},"passive",{get:function(){return A=!0,!0}});window.addEventListener("test",_,_),window.removeEventListener("test",_,_)}catch(e){A=!1}var T=!!A&&{passive:!1},O=function(e,t){var n=window.getComputedStyle(e);return"hidden"!==n[t]&&!(n.overflowY===n.overflowX&&"TEXTAREA"!==e.tagName&&"visible"===n[t])},x=function(e,t){var n=t;do{if("undefined"!=typeof ShadowRoot&&n instanceof ShadowRoot&&(n=n.host),N(e,n)){var r=I(e,n);if(r[1]>r[2])return!0}n=n.parentNode}while(n&&n!==document.body);return!1},N=function(e,t){return"v"===e?O(t,"overflowY"):O(t,"overflowX")},I=function(e,t){return"v"===e?[t.scrollTop,t.scrollHeight,t.clientHeight]:[t.scrollLeft,t.scrollWidth,t.clientWidth]},L=function(e,t,n,r,o){var a,c=(a=window.getComputedStyle(t).direction,"h"===e&&"rtl"===a?-1:1),u=c*r,i=n.target,l=t.contains(i),d=!1,s=u>0,f=0,p=0;do{var v=I(e,i),m=v[0],h=v[1]-v[2]-c*m;(m||h)&&N(e,i)&&(f+=h,p+=m),i=i.parentNode}while(!l&&i!==document.body||l&&(t.contains(i)||t===i));return s&&(o&&0===f||!o&&u>f)?d=!0:!s&&(o&&0===p||!o&&-u>p)&&(d=!0),d},P=function(e){return"changedTouches"in e?[e.changedTouches[0].clientX,e.changedTouches[0].clientY]:[0,0]},F=function(e){return[e.deltaX,e.deltaY]},W=function(e){return e&&"current"in e?e.current:e},Z=0,B=[],j=(v.useMedium(function(e){var t=s.useRef([]),n=s.useRef([0,0]),r=s.useRef(),o=s.useState(Z++)[0],a=s.useState(function(){return b()})[0],c=s.useRef(e);s.useEffect(function(){c.current=e},[e]),s.useEffect(function(){if(e.inert){document.body.classList.add("block-interactivity-".concat(o));var t=(function(e,t,n){if(n||2==arguments.length)for(var r,o=0,a=t.length;o<a;o++)!r&&o in t||(r||(r=Array.prototype.slice.call(t,0,o)),r[o]=t[o]);return e.concat(r||Array.prototype.slice.call(t))})([e.lockRef.current],(e.shards||[]).map(W),!0).filter(Boolean);return t.forEach(function(e){return e.classList.add("allow-interactivity-".concat(o))}),function(){document.body.classList.remove("block-interactivity-".concat(o)),t.forEach(function(e){return e.classList.remove("allow-interactivity-".concat(o))})}}},[e.inert,e.lockRef.current,e.shards]);var u=s.useCallback(function(e,t){if("touches"in e&&2===e.touches.length)return!c.current.allowPinchZoom;var o,a=P(e),u=n.current,i="deltaX"in e?e.deltaX:u[0]-a[0],l="deltaY"in e?e.deltaY:u[1]-a[1],d=e.target,s=Math.abs(i)>Math.abs(l)?"h":"v";if("touches"in e&&"h"===s&&"range"===d.type)return!1;var f=x(s,d);if(!f)return!0;if(f?o=s:(o="v"===s?"h":"v",f=x(s,d)),!f)return!1;if(!r.current&&"changedTouches"in e&&(i||l)&&(r.current=o),!o)return!0;var p=r.current||o;return L(p,t,e,"h"===p?i:l,!0)},[]),i=s.useCallback(function(e){if(B.length&&B[B.length-1]===a){var n="deltaY"in e?F(e):P(e),r=t.current.filter(function(t){var r;return t.name===e.type&&t.target===e.target&&(r=t.delta)[0]===n[0]&&r[1]===n[1]})[0];if(r&&r.should){e.cancelable&&e.preventDefault();return}if(!r){var o=(c.current.shards||[]).map(W).filter(Boolean).filter(function(t){return t.contains(e.target)});(o.length>0?u(e,o[0]):!c.current.noIsolation)&&e.cancelable&&e.preventDefault()}}},[]),l=s.useCallback(function(e,n,r,o){var a={name:e,delta:n,target:r,should:o};t.current.push(a),setTimeout(function(){t.current=t.current.filter(function(e){return e!==a})},1)},[]),d=s.useCallback(function(e){n.current=P(e),r.current=void 0},[]),f=s.useCallback(function(t){l(t.type,F(t),t.target,u(t,e.lockRef.current))},[]),p=s.useCallback(function(t){l(t.type,P(t),t.target,u(t,e.lockRef.current))},[]);s.useEffect(function(){return B.push(a),e.setCallbacks({onScrollCapture:f,onWheelCapture:f,onTouchMoveCapture:p}),document.addEventListener("wheel",i,T),document.addEventListener("touchmove",i,T),document.addEventListener("touchstart",d,T),function(){B=B.filter(function(e){return e!==a}),document.removeEventListener("wheel",i,T),document.removeEventListener("touchmove",i,T),document.removeEventListener("touchstart",d,T)}},[]);var v=e.removeScrollBar,m=e.inert;return s.createElement(s.Fragment,null,m?s.createElement(a,{styles:"\n  .block-interactivity-".concat(o," {pointer-events: none;}\n  .allow-interactivity-").concat(o," {pointer-events: all;}\n")}):null,v?s.createElement(k,{gapMode:"margin"}):null)}),g),K=s.forwardRef(function(e,t){return s.createElement(h,l({},e,{ref:t,sideCar:j}))});K.classNames=h.classNames;var V=K}}]);
(function(){const l=document.createElement("link").relList;if(l&&l.supports&&l.supports("modulepreload"))return;for(const f of document.querySelectorAll('link[rel="modulepreload"]'))u(f);new MutationObserver(f=>{for(const g of f)if(g.type==="childList")for(const h of g.addedNodes)h.tagName==="LINK"&&h.rel==="modulepreload"&&u(h)}).observe(document,{childList:!0,subtree:!0});function s(f){const g={};return f.integrity&&(g.integrity=f.integrity),f.referrerPolicy&&(g.referrerPolicy=f.referrerPolicy),f.crossOrigin==="use-credentials"?g.credentials="include":f.crossOrigin==="anonymous"?g.credentials="omit":g.credentials="same-origin",g}function u(f){if(f.ep)return;f.ep=!0;const g=s(f);fetch(f.href,g)}})();function Ut(n){return n&&n.__esModule&&Object.prototype.hasOwnProperty.call(n,"default")?n.default:n}var Ya={exports:{}},To={},qa={exports:{}},Ie={};/**
 * @license React
 * react.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var jd;function Zh(){if(jd)return Ie;jd=1;var n=Symbol.for("react.element"),l=Symbol.for("react.portal"),s=Symbol.for("react.fragment"),u=Symbol.for("react.strict_mode"),f=Symbol.for("react.profiler"),g=Symbol.for("react.provider"),h=Symbol.for("react.context"),v=Symbol.for("react.forward_ref"),x=Symbol.for("react.suspense"),E=Symbol.for("react.memo"),T=Symbol.for("react.lazy"),A=Symbol.iterator;function P(S){return S===null||typeof S!="object"?null:(S=A&&S[A]||S["@@iterator"],typeof S=="function"?S:null)}var U={isMounted:function(){return!1},enqueueForceUpdate:function(){},enqueueReplaceState:function(){},enqueueSetState:function(){}},j=Object.assign,y={};function $(S,C,I){this.props=S,this.context=C,this.refs=y,this.updater=I||U}$.prototype.isReactComponent={},$.prototype.setState=function(S,C){if(typeof S!="object"&&typeof S!="function"&&S!=null)throw Error("setState(...): takes an object of state variables to update or a function which returns an object of state variables.");this.updater.enqueueSetState(this,S,C,"setState")},$.prototype.forceUpdate=function(S){this.updater.enqueueForceUpdate(this,S,"forceUpdate")};function Y(){}Y.prototype=$.prototype;function H(S,C,I){this.props=S,this.context=C,this.refs=y,this.updater=I||U}var V=H.prototype=new Y;V.constructor=H,j(V,$.prototype),V.isPureReactComponent=!0;var J=Array.isArray,L=Object.prototype.hasOwnProperty,M={current:null},ne={key:!0,ref:!0,__self:!0,__source:!0};function de(S,C,I){var D,q={},b=null,oe=null;if(C!=null)for(D in C.ref!==void 0&&(oe=C.ref),C.key!==void 0&&(b=""+C.key),C)L.call(C,D)&&!ne.hasOwnProperty(D)&&(q[D]=C[D]);var Q=arguments.length-2;if(Q===1)q.children=I;else if(1<Q){for(var _=Array(Q),le=0;le<Q;le++)_[le]=arguments[le+2];q.children=_}if(S&&S.defaultProps)for(D in Q=S.defaultProps,Q)q[D]===void 0&&(q[D]=Q[D]);return{$$typeof:n,type:S,key:b,ref:oe,props:q,_owner:M.current}}function Ae(S,C){return{$$typeof:n,type:S.type,key:C,ref:S.ref,props:S.props,_owner:S._owner}}function $e(S){return typeof S=="object"&&S!==null&&S.$$typeof===n}function De(S){var C={"=":"=0",":":"=2"};return"$"+S.replace(/[=:]/g,function(I){return C[I]})}var W=/\/+/g;function ce(S,C){return typeof S=="object"&&S!==null&&S.key!=null?De(""+S.key):C.toString(36)}function xe(S,C,I,D,q){var b=typeof S;(b==="undefined"||b==="boolean")&&(S=null);var oe=!1;if(S===null)oe=!0;else switch(b){case"string":case"number":oe=!0;break;case"object":switch(S.$$typeof){case n:case l:oe=!0}}if(oe)return oe=S,q=q(oe),S=D===""?"."+ce(oe,0):D,J(q)?(I="",S!=null&&(I=S.replace(W,"$&/")+"/"),xe(q,C,I,"",function(le){return le})):q!=null&&($e(q)&&(q=Ae(q,I+(!q.key||oe&&oe.key===q.key?"":(""+q.key).replace(W,"$&/")+"/")+S)),C.push(q)),1;if(oe=0,D=D===""?".":D+":",J(S))for(var Q=0;Q<S.length;Q++){b=S[Q];var _=D+ce(b,Q);oe+=xe(b,C,I,_,q)}else if(_=P(S),typeof _=="function")for(S=_.call(S),Q=0;!(b=S.next()).done;)b=b.value,_=D+ce(b,Q++),oe+=xe(b,C,I,_,q);else if(b==="object")throw C=String(S),Error("Objects are not valid as a React child (found: "+(C==="[object Object]"?"object with keys {"+Object.keys(S).join(", ")+"}":C)+"). If you meant to render a collection of children, use an array instead.");return oe}function K(S,C,I){if(S==null)return S;var D=[],q=0;return xe(S,D,"","",function(b){return C.call(I,b,q++)}),D}function ue(S){if(S._status===-1){var C=S._result;C=C(),C.then(function(I){(S._status===0||S._status===-1)&&(S._status=1,S._result=I)},function(I){(S._status===0||S._status===-1)&&(S._status=2,S._result=I)}),S._status===-1&&(S._status=0,S._result=C)}if(S._status===1)return S._result.default;throw S._result}var ge={current:null},F={transition:null},se={ReactCurrentDispatcher:ge,ReactCurrentBatchConfig:F,ReactCurrentOwner:M};function G(){throw Error("act(...) is not supported in production builds of React.")}return Ie.Children={map:K,forEach:function(S,C,I){K(S,function(){C.apply(this,arguments)},I)},count:function(S){var C=0;return K(S,function(){C++}),C},toArray:function(S){return K(S,function(C){return C})||[]},only:function(S){if(!$e(S))throw Error("React.Children.only expected to receive a single React element child.");return S}},Ie.Component=$,Ie.Fragment=s,Ie.Profiler=f,Ie.PureComponent=H,Ie.StrictMode=u,Ie.Suspense=x,Ie.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=se,Ie.act=G,Ie.cloneElement=function(S,C,I){if(S==null)throw Error("React.cloneElement(...): The argument must be a React element, but you passed "+S+".");var D=j({},S.props),q=S.key,b=S.ref,oe=S._owner;if(C!=null){if(C.ref!==void 0&&(b=C.ref,oe=M.current),C.key!==void 0&&(q=""+C.key),S.type&&S.type.defaultProps)var Q=S.type.defaultProps;for(_ in C)L.call(C,_)&&!ne.hasOwnProperty(_)&&(D[_]=C[_]===void 0&&Q!==void 0?Q[_]:C[_])}var _=arguments.length-2;if(_===1)D.children=I;else if(1<_){Q=Array(_);for(var le=0;le<_;le++)Q[le]=arguments[le+2];D.children=Q}return{$$typeof:n,type:S.type,key:q,ref:b,props:D,_owner:oe}},Ie.createContext=function(S){return S={$$typeof:h,_currentValue:S,_currentValue2:S,_threadCount:0,Provider:null,Consumer:null,_defaultValue:null,_globalName:null},S.Provider={$$typeof:g,_context:S},S.Consumer=S},Ie.createElement=de,Ie.createFactory=function(S){var C=de.bind(null,S);return C.type=S,C},Ie.createRef=function(){return{current:null}},Ie.forwardRef=function(S){return{$$typeof:v,render:S}},Ie.isValidElement=$e,Ie.lazy=function(S){return{$$typeof:T,_payload:{_status:-1,_result:S},_init:ue}},Ie.memo=function(S,C){return{$$typeof:E,type:S,compare:C===void 0?null:C}},Ie.startTransition=function(S){var C=F.transition;F.transition={};try{S()}finally{F.transition=C}},Ie.unstable_act=G,Ie.useCallback=function(S,C){return ge.current.useCallback(S,C)},Ie.useContext=function(S){return ge.current.useContext(S)},Ie.useDebugValue=function(){},Ie.useDeferredValue=function(S){return ge.current.useDeferredValue(S)},Ie.useEffect=function(S,C){return ge.current.useEffect(S,C)},Ie.useId=function(){return ge.current.useId()},Ie.useImperativeHandle=function(S,C,I){return ge.current.useImperativeHandle(S,C,I)},Ie.useInsertionEffect=function(S,C){return ge.current.useInsertionEffect(S,C)},Ie.useLayoutEffect=function(S,C){return ge.current.useLayoutEffect(S,C)},Ie.useMemo=function(S,C){return ge.current.useMemo(S,C)},Ie.useReducer=function(S,C,I){return ge.current.useReducer(S,C,I)},Ie.useRef=function(S){return ge.current.useRef(S)},Ie.useState=function(S){return ge.current.useState(S)},Ie.useSyncExternalStore=function(S,C,I){return ge.current.useSyncExternalStore(S,C,I)},Ie.useTransition=function(){return ge.current.useTransition()},Ie.version="18.3.1",Ie}var Pd;function Rs(){return Pd||(Pd=1,qa.exports=Zh()),qa.exports}/**
 * @license React
 * react-jsx-runtime.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var Rd;function Jh(){if(Rd)return To;Rd=1;var n=Rs(),l=Symbol.for("react.element"),s=Symbol.for("react.fragment"),u=Object.prototype.hasOwnProperty,f=n.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED.ReactCurrentOwner,g={key:!0,ref:!0,__self:!0,__source:!0};function h(v,x,E){var T,A={},P=null,U=null;E!==void 0&&(P=""+E),x.key!==void 0&&(P=""+x.key),x.ref!==void 0&&(U=x.ref);for(T in x)u.call(x,T)&&!g.hasOwnProperty(T)&&(A[T]=x[T]);if(v&&v.defaultProps)for(T in x=v.defaultProps,x)A[T]===void 0&&(A[T]=x[T]);return{$$typeof:l,type:v,key:P,ref:U,props:A,_owner:f.current}}return To.Fragment=s,To.jsx=h,To.jsxs=h,To}var $d;function em(){return $d||($d=1,Ya.exports=Jh()),Ya.exports}var c=em(),m=Rs();const B=Ut(m);var Ui={},Xa={exports:{}},Pt={},Za={exports:{}},Ja={};/**
 * @license React
 * scheduler.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var Dd;function tm(){return Dd||(Dd=1,(function(n){function l(F,se){var G=F.length;F.push(se);e:for(;0<G;){var S=G-1>>>1,C=F[S];if(0<f(C,se))F[S]=se,F[G]=C,G=S;else break e}}function s(F){return F.length===0?null:F[0]}function u(F){if(F.length===0)return null;var se=F[0],G=F.pop();if(G!==se){F[0]=G;e:for(var S=0,C=F.length,I=C>>>1;S<I;){var D=2*(S+1)-1,q=F[D],b=D+1,oe=F[b];if(0>f(q,G))b<C&&0>f(oe,q)?(F[S]=oe,F[b]=G,S=b):(F[S]=q,F[D]=G,S=D);else if(b<C&&0>f(oe,G))F[S]=oe,F[b]=G,S=b;else break e}}return se}function f(F,se){var G=F.sortIndex-se.sortIndex;return G!==0?G:F.id-se.id}if(typeof performance=="object"&&typeof performance.now=="function"){var g=performance;n.unstable_now=function(){return g.now()}}else{var h=Date,v=h.now();n.unstable_now=function(){return h.now()-v}}var x=[],E=[],T=1,A=null,P=3,U=!1,j=!1,y=!1,$=typeof setTimeout=="function"?setTimeout:null,Y=typeof clearTimeout=="function"?clearTimeout:null,H=typeof setImmediate<"u"?setImmediate:null;typeof navigator<"u"&&navigator.scheduling!==void 0&&navigator.scheduling.isInputPending!==void 0&&navigator.scheduling.isInputPending.bind(navigator.scheduling);function V(F){for(var se=s(E);se!==null;){if(se.callback===null)u(E);else if(se.startTime<=F)u(E),se.sortIndex=se.expirationTime,l(x,se);else break;se=s(E)}}function J(F){if(y=!1,V(F),!j)if(s(x)!==null)j=!0,ue(L);else{var se=s(E);se!==null&&ge(J,se.startTime-F)}}function L(F,se){j=!1,y&&(y=!1,Y(de),de=-1),U=!0;var G=P;try{for(V(se),A=s(x);A!==null&&(!(A.expirationTime>se)||F&&!De());){var S=A.callback;if(typeof S=="function"){A.callback=null,P=A.priorityLevel;var C=S(A.expirationTime<=se);se=n.unstable_now(),typeof C=="function"?A.callback=C:A===s(x)&&u(x),V(se)}else u(x);A=s(x)}if(A!==null)var I=!0;else{var D=s(E);D!==null&&ge(J,D.startTime-se),I=!1}return I}finally{A=null,P=G,U=!1}}var M=!1,ne=null,de=-1,Ae=5,$e=-1;function De(){return!(n.unstable_now()-$e<Ae)}function W(){if(ne!==null){var F=n.unstable_now();$e=F;var se=!0;try{se=ne(!0,F)}finally{se?ce():(M=!1,ne=null)}}else M=!1}var ce;if(typeof H=="function")ce=function(){H(W)};else if(typeof MessageChannel<"u"){var xe=new MessageChannel,K=xe.port2;xe.port1.onmessage=W,ce=function(){K.postMessage(null)}}else ce=function(){$(W,0)};function ue(F){ne=F,M||(M=!0,ce())}function ge(F,se){de=$(function(){F(n.unstable_now())},se)}n.unstable_IdlePriority=5,n.unstable_ImmediatePriority=1,n.unstable_LowPriority=4,n.unstable_NormalPriority=3,n.unstable_Profiling=null,n.unstable_UserBlockingPriority=2,n.unstable_cancelCallback=function(F){F.callback=null},n.unstable_continueExecution=function(){j||U||(j=!0,ue(L))},n.unstable_forceFrameRate=function(F){0>F||125<F?console.error("forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported"):Ae=0<F?Math.floor(1e3/F):5},n.unstable_getCurrentPriorityLevel=function(){return P},n.unstable_getFirstCallbackNode=function(){return s(x)},n.unstable_next=function(F){switch(P){case 1:case 2:case 3:var se=3;break;default:se=P}var G=P;P=se;try{return F()}finally{P=G}},n.unstable_pauseExecution=function(){},n.unstable_requestPaint=function(){},n.unstable_runWithPriority=function(F,se){switch(F){case 1:case 2:case 3:case 4:case 5:break;default:F=3}var G=P;P=F;try{return se()}finally{P=G}},n.unstable_scheduleCallback=function(F,se,G){var S=n.unstable_now();switch(typeof G=="object"&&G!==null?(G=G.delay,G=typeof G=="number"&&0<G?S+G:S):G=S,F){case 1:var C=-1;break;case 2:C=250;break;case 5:C=1073741823;break;case 4:C=1e4;break;default:C=5e3}return C=G+C,F={id:T++,callback:se,priorityLevel:F,startTime:G,expirationTime:C,sortIndex:-1},G>S?(F.sortIndex=G,l(E,F),s(x)===null&&F===s(E)&&(y?(Y(de),de=-1):y=!0,ge(J,G-S))):(F.sortIndex=C,l(x,F),j||U||(j=!0,ue(L))),F},n.unstable_shouldYield=De,n.unstable_wrapCallback=function(F){var se=P;return function(){var G=P;P=se;try{return F.apply(this,arguments)}finally{P=G}}}})(Ja)),Ja}var zd;function rm(){return zd||(zd=1,Za.exports=tm()),Za.exports}/**
 * @license React
 * react-dom.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var Id;function nm(){if(Id)return Pt;Id=1;var n=Rs(),l=rm();function s(e){for(var t="https://reactjs.org/docs/error-decoder.html?invariant="+e,r=1;r<arguments.length;r++)t+="&args[]="+encodeURIComponent(arguments[r]);return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}var u=new Set,f={};function g(e,t){h(e,t),h(e+"Capture",t)}function h(e,t){for(f[e]=t,e=0;e<t.length;e++)u.add(t[e])}var v=!(typeof window>"u"||typeof window.document>"u"||typeof window.document.createElement>"u"),x=Object.prototype.hasOwnProperty,E=/^[:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD][:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD\-.0-9\u00B7\u0300-\u036F\u203F-\u2040]*$/,T={},A={};function P(e){return x.call(A,e)?!0:x.call(T,e)?!1:E.test(e)?A[e]=!0:(T[e]=!0,!1)}function U(e,t,r,o){if(r!==null&&r.type===0)return!1;switch(typeof t){case"function":case"symbol":return!0;case"boolean":return o?!1:r!==null?!r.acceptsBooleans:(e=e.toLowerCase().slice(0,5),e!=="data-"&&e!=="aria-");default:return!1}}function j(e,t,r,o){if(t===null||typeof t>"u"||U(e,t,r,o))return!0;if(o)return!1;if(r!==null)switch(r.type){case 3:return!t;case 4:return t===!1;case 5:return isNaN(t);case 6:return isNaN(t)||1>t}return!1}function y(e,t,r,o,i,a,d){this.acceptsBooleans=t===2||t===3||t===4,this.attributeName=o,this.attributeNamespace=i,this.mustUseProperty=r,this.propertyName=e,this.type=t,this.sanitizeURL=a,this.removeEmptyString=d}var $={};"children dangerouslySetInnerHTML defaultValue defaultChecked innerHTML suppressContentEditableWarning suppressHydrationWarning style".split(" ").forEach(function(e){$[e]=new y(e,0,!1,e,null,!1,!1)}),[["acceptCharset","accept-charset"],["className","class"],["htmlFor","for"],["httpEquiv","http-equiv"]].forEach(function(e){var t=e[0];$[t]=new y(t,1,!1,e[1],null,!1,!1)}),["contentEditable","draggable","spellCheck","value"].forEach(function(e){$[e]=new y(e,2,!1,e.toLowerCase(),null,!1,!1)}),["autoReverse","externalResourcesRequired","focusable","preserveAlpha"].forEach(function(e){$[e]=new y(e,2,!1,e,null,!1,!1)}),"allowFullScreen async autoFocus autoPlay controls default defer disabled disablePictureInPicture disableRemotePlayback formNoValidate hidden loop noModule noValidate open playsInline readOnly required reversed scoped seamless itemScope".split(" ").forEach(function(e){$[e]=new y(e,3,!1,e.toLowerCase(),null,!1,!1)}),["checked","multiple","muted","selected"].forEach(function(e){$[e]=new y(e,3,!0,e,null,!1,!1)}),["capture","download"].forEach(function(e){$[e]=new y(e,4,!1,e,null,!1,!1)}),["cols","rows","size","span"].forEach(function(e){$[e]=new y(e,6,!1,e,null,!1,!1)}),["rowSpan","start"].forEach(function(e){$[e]=new y(e,5,!1,e.toLowerCase(),null,!1,!1)});var Y=/[\-:]([a-z])/g;function H(e){return e[1].toUpperCase()}"accent-height alignment-baseline arabic-form baseline-shift cap-height clip-path clip-rule color-interpolation color-interpolation-filters color-profile color-rendering dominant-baseline enable-background fill-opacity fill-rule flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight glyph-name glyph-orientation-horizontal glyph-orientation-vertical horiz-adv-x horiz-origin-x image-rendering letter-spacing lighting-color marker-end marker-mid marker-start overline-position overline-thickness paint-order panose-1 pointer-events rendering-intent shape-rendering stop-color stop-opacity strikethrough-position strikethrough-thickness stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width text-anchor text-decoration text-rendering underline-position underline-thickness unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical vector-effect vert-adv-y vert-origin-x vert-origin-y word-spacing writing-mode xmlns:xlink x-height".split(" ").forEach(function(e){var t=e.replace(Y,H);$[t]=new y(t,1,!1,e,null,!1,!1)}),"xlink:actuate xlink:arcrole xlink:role xlink:show xlink:title xlink:type".split(" ").forEach(function(e){var t=e.replace(Y,H);$[t]=new y(t,1,!1,e,"http://www.w3.org/1999/xlink",!1,!1)}),["xml:base","xml:lang","xml:space"].forEach(function(e){var t=e.replace(Y,H);$[t]=new y(t,1,!1,e,"http://www.w3.org/XML/1998/namespace",!1,!1)}),["tabIndex","crossOrigin"].forEach(function(e){$[e]=new y(e,1,!1,e.toLowerCase(),null,!1,!1)}),$.xlinkHref=new y("xlinkHref",1,!1,"xlink:href","http://www.w3.org/1999/xlink",!0,!1),["src","href","action","formAction"].forEach(function(e){$[e]=new y(e,1,!1,e.toLowerCase(),null,!0,!0)});function V(e,t,r,o){var i=$.hasOwnProperty(t)?$[t]:null;(i!==null?i.type!==0:o||!(2<t.length)||t[0]!=="o"&&t[0]!=="O"||t[1]!=="n"&&t[1]!=="N")&&(j(t,r,i,o)&&(r=null),o||i===null?P(t)&&(r===null?e.removeAttribute(t):e.setAttribute(t,""+r)):i.mustUseProperty?e[i.propertyName]=r===null?i.type===3?!1:"":r:(t=i.attributeName,o=i.attributeNamespace,r===null?e.removeAttribute(t):(i=i.type,r=i===3||i===4&&r===!0?"":""+r,o?e.setAttributeNS(o,t,r):e.setAttribute(t,r))))}var J=n.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED,L=Symbol.for("react.element"),M=Symbol.for("react.portal"),ne=Symbol.for("react.fragment"),de=Symbol.for("react.strict_mode"),Ae=Symbol.for("react.profiler"),$e=Symbol.for("react.provider"),De=Symbol.for("react.context"),W=Symbol.for("react.forward_ref"),ce=Symbol.for("react.suspense"),xe=Symbol.for("react.suspense_list"),K=Symbol.for("react.memo"),ue=Symbol.for("react.lazy"),ge=Symbol.for("react.offscreen"),F=Symbol.iterator;function se(e){return e===null||typeof e!="object"?null:(e=F&&e[F]||e["@@iterator"],typeof e=="function"?e:null)}var G=Object.assign,S;function C(e){if(S===void 0)try{throw Error()}catch(r){var t=r.stack.trim().match(/\n( *(at )?)/);S=t&&t[1]||""}return`
`+S+e}var I=!1;function D(e,t){if(!e||I)return"";I=!0;var r=Error.prepareStackTrace;Error.prepareStackTrace=void 0;try{if(t)if(t=function(){throw Error()},Object.defineProperty(t.prototype,"props",{set:function(){throw Error()}}),typeof Reflect=="object"&&Reflect.construct){try{Reflect.construct(t,[])}catch(O){var o=O}Reflect.construct(e,[],t)}else{try{t.call()}catch(O){o=O}e.call(t.prototype)}else{try{throw Error()}catch(O){o=O}e()}}catch(O){if(O&&o&&typeof O.stack=="string"){for(var i=O.stack.split(`
`),a=o.stack.split(`
`),d=i.length-1,p=a.length-1;1<=d&&0<=p&&i[d]!==a[p];)p--;for(;1<=d&&0<=p;d--,p--)if(i[d]!==a[p]){if(d!==1||p!==1)do if(d--,p--,0>p||i[d]!==a[p]){var w=`
`+i[d].replace(" at new "," at ");return e.displayName&&w.includes("<anonymous>")&&(w=w.replace("<anonymous>",e.displayName)),w}while(1<=d&&0<=p);break}}}finally{I=!1,Error.prepareStackTrace=r}return(e=e?e.displayName||e.name:"")?C(e):""}function q(e){switch(e.tag){case 5:return C(e.type);case 16:return C("Lazy");case 13:return C("Suspense");case 19:return C("SuspenseList");case 0:case 2:case 15:return e=D(e.type,!1),e;case 11:return e=D(e.type.render,!1),e;case 1:return e=D(e.type,!0),e;default:return""}}function b(e){if(e==null)return null;if(typeof e=="function")return e.displayName||e.name||null;if(typeof e=="string")return e;switch(e){case ne:return"Fragment";case M:return"Portal";case Ae:return"Profiler";case de:return"StrictMode";case ce:return"Suspense";case xe:return"SuspenseList"}if(typeof e=="object")switch(e.$$typeof){case De:return(e.displayName||"Context")+".Consumer";case $e:return(e._context.displayName||"Context")+".Provider";case W:var t=e.render;return e=e.displayName,e||(e=t.displayName||t.name||"",e=e!==""?"ForwardRef("+e+")":"ForwardRef"),e;case K:return t=e.displayName||null,t!==null?t:b(e.type)||"Memo";case ue:t=e._payload,e=e._init;try{return b(e(t))}catch{}}return null}function oe(e){var t=e.type;switch(e.tag){case 24:return"Cache";case 9:return(t.displayName||"Context")+".Consumer";case 10:return(t._context.displayName||"Context")+".Provider";case 18:return"DehydratedFragment";case 11:return e=t.render,e=e.displayName||e.name||"",t.displayName||(e!==""?"ForwardRef("+e+")":"ForwardRef");case 7:return"Fragment";case 5:return t;case 4:return"Portal";case 3:return"Root";case 6:return"Text";case 16:return b(t);case 8:return t===de?"StrictMode":"Mode";case 22:return"Offscreen";case 12:return"Profiler";case 21:return"Scope";case 13:return"Suspense";case 19:return"SuspenseList";case 25:return"TracingMarker";case 1:case 0:case 17:case 2:case 14:case 15:if(typeof t=="function")return t.displayName||t.name||null;if(typeof t=="string")return t}return null}function Q(e){switch(typeof e){case"boolean":case"number":case"string":case"undefined":return e;case"object":return e;default:return""}}function _(e){var t=e.type;return(e=e.nodeName)&&e.toLowerCase()==="input"&&(t==="checkbox"||t==="radio")}function le(e){var t=_(e)?"checked":"value",r=Object.getOwnPropertyDescriptor(e.constructor.prototype,t),o=""+e[t];if(!e.hasOwnProperty(t)&&typeof r<"u"&&typeof r.get=="function"&&typeof r.set=="function"){var i=r.get,a=r.set;return Object.defineProperty(e,t,{configurable:!0,get:function(){return i.call(this)},set:function(d){o=""+d,a.call(this,d)}}),Object.defineProperty(e,t,{enumerable:r.enumerable}),{getValue:function(){return o},setValue:function(d){o=""+d},stopTracking:function(){e._valueTracker=null,delete e[t]}}}}function fe(e){e._valueTracker||(e._valueTracker=le(e))}function ae(e){if(!e)return!1;var t=e._valueTracker;if(!t)return!0;var r=t.getValue(),o="";return e&&(o=_(e)?e.checked?"true":"false":e.value),e=o,e!==r?(t.setValue(e),!0):!1}function je(e){if(e=e||(typeof document<"u"?document:void 0),typeof e>"u")return null;try{return e.activeElement||e.body}catch{return e.body}}function ze(e,t){var r=t.checked;return G({},t,{defaultChecked:void 0,defaultValue:void 0,value:void 0,checked:r??e._wrapperState.initialChecked})}function Ce(e,t){var r=t.defaultValue==null?"":t.defaultValue,o=t.checked!=null?t.checked:t.defaultChecked;r=Q(t.value!=null?t.value:r),e._wrapperState={initialChecked:o,initialValue:r,controlled:t.type==="checkbox"||t.type==="radio"?t.checked!=null:t.value!=null}}function Me(e,t){t=t.checked,t!=null&&V(e,"checked",t,!1)}function Re(e,t){Me(e,t);var r=Q(t.value),o=t.type;if(r!=null)o==="number"?(r===0&&e.value===""||e.value!=r)&&(e.value=""+r):e.value!==""+r&&(e.value=""+r);else if(o==="submit"||o==="reset"){e.removeAttribute("value");return}t.hasOwnProperty("value")?ot(e,t.type,r):t.hasOwnProperty("defaultValue")&&ot(e,t.type,Q(t.defaultValue)),t.checked==null&&t.defaultChecked!=null&&(e.defaultChecked=!!t.defaultChecked)}function Ne(e,t,r){if(t.hasOwnProperty("value")||t.hasOwnProperty("defaultValue")){var o=t.type;if(!(o!=="submit"&&o!=="reset"||t.value!==void 0&&t.value!==null))return;t=""+e._wrapperState.initialValue,r||t===e.value||(e.value=t),e.defaultValue=t}r=e.name,r!==""&&(e.name=""),e.defaultChecked=!!e._wrapperState.initialChecked,r!==""&&(e.name=r)}function ot(e,t,r){(t!=="number"||je(e.ownerDocument)!==e)&&(r==null?e.defaultValue=""+e._wrapperState.initialValue:e.defaultValue!==""+r&&(e.defaultValue=""+r))}var Fe=Array.isArray;function N(e,t,r,o){if(e=e.options,t){t={};for(var i=0;i<r.length;i++)t["$"+r[i]]=!0;for(r=0;r<e.length;r++)i=t.hasOwnProperty("$"+e[r].value),e[r].selected!==i&&(e[r].selected=i),i&&o&&(e[r].defaultSelected=!0)}else{for(r=""+Q(r),t=null,i=0;i<e.length;i++){if(e[i].value===r){e[i].selected=!0,o&&(e[i].defaultSelected=!0);return}t!==null||e[i].disabled||(t=e[i])}t!==null&&(t.selected=!0)}}function te(e,t){if(t.dangerouslySetInnerHTML!=null)throw Error(s(91));return G({},t,{value:void 0,defaultValue:void 0,children:""+e._wrapperState.initialValue})}function Ee(e,t){var r=t.value;if(r==null){if(r=t.children,t=t.defaultValue,r!=null){if(t!=null)throw Error(s(92));if(Fe(r)){if(1<r.length)throw Error(s(93));r=r[0]}t=r}t==null&&(t=""),r=t}e._wrapperState={initialValue:Q(r)}}function ve(e,t){var r=Q(t.value),o=Q(t.defaultValue);r!=null&&(r=""+r,r!==e.value&&(e.value=r),t.defaultValue==null&&e.defaultValue!==r&&(e.defaultValue=r)),o!=null&&(e.defaultValue=""+o)}function Ue(e){var t=e.textContent;t===e._wrapperState.initialValue&&t!==""&&t!==null&&(e.value=t)}function ur(e){switch(e){case"svg":return"http://www.w3.org/2000/svg";case"math":return"http://www.w3.org/1998/Math/MathML";default:return"http://www.w3.org/1999/xhtml"}}function Hr(e,t){return e==null||e==="http://www.w3.org/1999/xhtml"?ur(t):e==="http://www.w3.org/2000/svg"&&t==="foreignObject"?"http://www.w3.org/1999/xhtml":e}var Rt,Hs=(function(e){return typeof MSApp<"u"&&MSApp.execUnsafeLocalFunction?function(t,r,o,i){MSApp.execUnsafeLocalFunction(function(){return e(t,r,o,i)})}:e})(function(e,t){if(e.namespaceURI!=="http://www.w3.org/2000/svg"||"innerHTML"in e)e.innerHTML=t;else{for(Rt=Rt||document.createElement("div"),Rt.innerHTML="<svg>"+t.valueOf().toString()+"</svg>",t=Rt.firstChild;e.firstChild;)e.removeChild(e.firstChild);for(;t.firstChild;)e.appendChild(t.firstChild)}});function Hn(e,t){if(t){var r=e.firstChild;if(r&&r===e.lastChild&&r.nodeType===3){r.nodeValue=t;return}}e.textContent=t}var Un={animationIterationCount:!0,aspectRatio:!0,borderImageOutset:!0,borderImageSlice:!0,borderImageWidth:!0,boxFlex:!0,boxFlexGroup:!0,boxOrdinalGroup:!0,columnCount:!0,columns:!0,flex:!0,flexGrow:!0,flexPositive:!0,flexShrink:!0,flexNegative:!0,flexOrder:!0,gridArea:!0,gridRow:!0,gridRowEnd:!0,gridRowSpan:!0,gridRowStart:!0,gridColumn:!0,gridColumnEnd:!0,gridColumnSpan:!0,gridColumnStart:!0,fontWeight:!0,lineClamp:!0,lineHeight:!0,opacity:!0,order:!0,orphans:!0,tabSize:!0,widows:!0,zIndex:!0,zoom:!0,fillOpacity:!0,floodOpacity:!0,stopOpacity:!0,strokeDasharray:!0,strokeDashoffset:!0,strokeMiterlimit:!0,strokeOpacity:!0,strokeWidth:!0},rp=["Webkit","ms","Moz","O"];Object.keys(Un).forEach(function(e){rp.forEach(function(t){t=t+e.charAt(0).toUpperCase()+e.substring(1),Un[t]=Un[e]})});function Us(e,t,r){return t==null||typeof t=="boolean"||t===""?"":r||typeof t!="number"||t===0||Un.hasOwnProperty(e)&&Un[e]?(""+t).trim():t+"px"}function Vs(e,t){e=e.style;for(var r in t)if(t.hasOwnProperty(r)){var o=r.indexOf("--")===0,i=Us(r,t[r],o);r==="float"&&(r="cssFloat"),o?e.setProperty(r,i):e[r]=i}}var np=G({menuitem:!0},{area:!0,base:!0,br:!0,col:!0,embed:!0,hr:!0,img:!0,input:!0,keygen:!0,link:!0,meta:!0,param:!0,source:!0,track:!0,wbr:!0});function al(e,t){if(t){if(np[e]&&(t.children!=null||t.dangerouslySetInnerHTML!=null))throw Error(s(137,e));if(t.dangerouslySetInnerHTML!=null){if(t.children!=null)throw Error(s(60));if(typeof t.dangerouslySetInnerHTML!="object"||!("__html"in t.dangerouslySetInnerHTML))throw Error(s(61))}if(t.style!=null&&typeof t.style!="object")throw Error(s(62))}}function sl(e,t){if(e.indexOf("-")===-1)return typeof t.is=="string";switch(e){case"annotation-xml":case"color-profile":case"font-face":case"font-face-src":case"font-face-uri":case"font-face-format":case"font-face-name":case"missing-glyph":return!1;default:return!0}}var cl=null;function ul(e){return e=e.target||e.srcElement||window,e.correspondingUseElement&&(e=e.correspondingUseElement),e.nodeType===3?e.parentNode:e}var dl=null,fn=null,pn=null;function Qs(e){if(e=po(e)){if(typeof dl!="function")throw Error(s(280));var t=e.stateNode;t&&(t=ii(t),dl(e.stateNode,e.type,t))}}function Gs(e){fn?pn?pn.push(e):pn=[e]:fn=e}function Ks(){if(fn){var e=fn,t=pn;if(pn=fn=null,Qs(e),t)for(e=0;e<t.length;e++)Qs(t[e])}}function Ys(e,t){return e(t)}function qs(){}var fl=!1;function Xs(e,t,r){if(fl)return e(t,r);fl=!0;try{return Ys(e,t,r)}finally{fl=!1,(fn!==null||pn!==null)&&(qs(),Ks())}}function Vn(e,t){var r=e.stateNode;if(r===null)return null;var o=ii(r);if(o===null)return null;r=o[t];e:switch(t){case"onClick":case"onClickCapture":case"onDoubleClick":case"onDoubleClickCapture":case"onMouseDown":case"onMouseDownCapture":case"onMouseMove":case"onMouseMoveCapture":case"onMouseUp":case"onMouseUpCapture":case"onMouseEnter":(o=!o.disabled)||(e=e.type,o=!(e==="button"||e==="input"||e==="select"||e==="textarea")),e=!o;break e;default:e=!1}if(e)return null;if(r&&typeof r!="function")throw Error(s(231,t,typeof r));return r}var pl=!1;if(v)try{var Qn={};Object.defineProperty(Qn,"passive",{get:function(){pl=!0}}),window.addEventListener("test",Qn,Qn),window.removeEventListener("test",Qn,Qn)}catch{pl=!1}function op(e,t,r,o,i,a,d,p,w){var O=Array.prototype.slice.call(arguments,3);try{t.apply(r,O)}catch(Z){this.onError(Z)}}var Gn=!1,Mo=null,Oo=!1,hl=null,ip={onError:function(e){Gn=!0,Mo=e}};function lp(e,t,r,o,i,a,d,p,w){Gn=!1,Mo=null,op.apply(ip,arguments)}function ap(e,t,r,o,i,a,d,p,w){if(lp.apply(this,arguments),Gn){if(Gn){var O=Mo;Gn=!1,Mo=null}else throw Error(s(198));Oo||(Oo=!0,hl=O)}}function Ur(e){var t=e,r=e;if(e.alternate)for(;t.return;)t=t.return;else{e=t;do t=e,(t.flags&4098)!==0&&(r=t.return),e=t.return;while(e)}return t.tag===3?r:null}function Zs(e){if(e.tag===13){var t=e.memoizedState;if(t===null&&(e=e.alternate,e!==null&&(t=e.memoizedState)),t!==null)return t.dehydrated}return null}function Js(e){if(Ur(e)!==e)throw Error(s(188))}function sp(e){var t=e.alternate;if(!t){if(t=Ur(e),t===null)throw Error(s(188));return t!==e?null:e}for(var r=e,o=t;;){var i=r.return;if(i===null)break;var a=i.alternate;if(a===null){if(o=i.return,o!==null){r=o;continue}break}if(i.child===a.child){for(a=i.child;a;){if(a===r)return Js(i),e;if(a===o)return Js(i),t;a=a.sibling}throw Error(s(188))}if(r.return!==o.return)r=i,o=a;else{for(var d=!1,p=i.child;p;){if(p===r){d=!0,r=i,o=a;break}if(p===o){d=!0,o=i,r=a;break}p=p.sibling}if(!d){for(p=a.child;p;){if(p===r){d=!0,r=a,o=i;break}if(p===o){d=!0,o=a,r=i;break}p=p.sibling}if(!d)throw Error(s(189))}}if(r.alternate!==o)throw Error(s(190))}if(r.tag!==3)throw Error(s(188));return r.stateNode.current===r?e:t}function ec(e){return e=sp(e),e!==null?tc(e):null}function tc(e){if(e.tag===5||e.tag===6)return e;for(e=e.child;e!==null;){var t=tc(e);if(t!==null)return t;e=e.sibling}return null}var rc=l.unstable_scheduleCallback,nc=l.unstable_cancelCallback,cp=l.unstable_shouldYield,up=l.unstable_requestPaint,it=l.unstable_now,dp=l.unstable_getCurrentPriorityLevel,ml=l.unstable_ImmediatePriority,oc=l.unstable_UserBlockingPriority,No=l.unstable_NormalPriority,fp=l.unstable_LowPriority,ic=l.unstable_IdlePriority,Bo=null,tr=null;function pp(e){if(tr&&typeof tr.onCommitFiberRoot=="function")try{tr.onCommitFiberRoot(Bo,e,void 0,(e.current.flags&128)===128)}catch{}}var Vt=Math.clz32?Math.clz32:gp,hp=Math.log,mp=Math.LN2;function gp(e){return e>>>=0,e===0?32:31-(hp(e)/mp|0)|0}var Fo=64,Wo=4194304;function Kn(e){switch(e&-e){case 1:return 1;case 2:return 2;case 4:return 4;case 8:return 8;case 16:return 16;case 32:return 32;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return e&4194240;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return e&130023424;case 134217728:return 134217728;case 268435456:return 268435456;case 536870912:return 536870912;case 1073741824:return 1073741824;default:return e}}function Ho(e,t){var r=e.pendingLanes;if(r===0)return 0;var o=0,i=e.suspendedLanes,a=e.pingedLanes,d=r&268435455;if(d!==0){var p=d&~i;p!==0?o=Kn(p):(a&=d,a!==0&&(o=Kn(a)))}else d=r&~i,d!==0?o=Kn(d):a!==0&&(o=Kn(a));if(o===0)return 0;if(t!==0&&t!==o&&(t&i)===0&&(i=o&-o,a=t&-t,i>=a||i===16&&(a&4194240)!==0))return t;if((o&4)!==0&&(o|=r&16),t=e.entangledLanes,t!==0)for(e=e.entanglements,t&=o;0<t;)r=31-Vt(t),i=1<<r,o|=e[r],t&=~i;return o}function vp(e,t){switch(e){case 1:case 2:case 4:return t+250;case 8:case 16:case 32:case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return t+5e3;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return-1;case 134217728:case 268435456:case 536870912:case 1073741824:return-1;default:return-1}}function yp(e,t){for(var r=e.suspendedLanes,o=e.pingedLanes,i=e.expirationTimes,a=e.pendingLanes;0<a;){var d=31-Vt(a),p=1<<d,w=i[d];w===-1?((p&r)===0||(p&o)!==0)&&(i[d]=vp(p,t)):w<=t&&(e.expiredLanes|=p),a&=~p}}function gl(e){return e=e.pendingLanes&-1073741825,e!==0?e:e&1073741824?1073741824:0}function lc(){var e=Fo;return Fo<<=1,(Fo&4194240)===0&&(Fo=64),e}function vl(e){for(var t=[],r=0;31>r;r++)t.push(e);return t}function Yn(e,t,r){e.pendingLanes|=t,t!==536870912&&(e.suspendedLanes=0,e.pingedLanes=0),e=e.eventTimes,t=31-Vt(t),e[t]=r}function xp(e,t){var r=e.pendingLanes&~t;e.pendingLanes=t,e.suspendedLanes=0,e.pingedLanes=0,e.expiredLanes&=t,e.mutableReadLanes&=t,e.entangledLanes&=t,t=e.entanglements;var o=e.eventTimes;for(e=e.expirationTimes;0<r;){var i=31-Vt(r),a=1<<i;t[i]=0,o[i]=-1,e[i]=-1,r&=~a}}function yl(e,t){var r=e.entangledLanes|=t;for(e=e.entanglements;r;){var o=31-Vt(r),i=1<<o;i&t|e[o]&t&&(e[o]|=t),r&=~i}}var Ve=0;function ac(e){return e&=-e,1<e?4<e?(e&268435455)!==0?16:536870912:4:1}var sc,xl,cc,uc,dc,wl=!1,Uo=[],br=null,kr=null,Sr=null,qn=new Map,Xn=new Map,Cr=[],wp="mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset submit".split(" ");function fc(e,t){switch(e){case"focusin":case"focusout":br=null;break;case"dragenter":case"dragleave":kr=null;break;case"mouseover":case"mouseout":Sr=null;break;case"pointerover":case"pointerout":qn.delete(t.pointerId);break;case"gotpointercapture":case"lostpointercapture":Xn.delete(t.pointerId)}}function Zn(e,t,r,o,i,a){return e===null||e.nativeEvent!==a?(e={blockedOn:t,domEventName:r,eventSystemFlags:o,nativeEvent:a,targetContainers:[i]},t!==null&&(t=po(t),t!==null&&xl(t)),e):(e.eventSystemFlags|=o,t=e.targetContainers,i!==null&&t.indexOf(i)===-1&&t.push(i),e)}function bp(e,t,r,o,i){switch(t){case"focusin":return br=Zn(br,e,t,r,o,i),!0;case"dragenter":return kr=Zn(kr,e,t,r,o,i),!0;case"mouseover":return Sr=Zn(Sr,e,t,r,o,i),!0;case"pointerover":var a=i.pointerId;return qn.set(a,Zn(qn.get(a)||null,e,t,r,o,i)),!0;case"gotpointercapture":return a=i.pointerId,Xn.set(a,Zn(Xn.get(a)||null,e,t,r,o,i)),!0}return!1}function pc(e){var t=Vr(e.target);if(t!==null){var r=Ur(t);if(r!==null){if(t=r.tag,t===13){if(t=Zs(r),t!==null){e.blockedOn=t,dc(e.priority,function(){cc(r)});return}}else if(t===3&&r.stateNode.current.memoizedState.isDehydrated){e.blockedOn=r.tag===3?r.stateNode.containerInfo:null;return}}}e.blockedOn=null}function Vo(e){if(e.blockedOn!==null)return!1;for(var t=e.targetContainers;0<t.length;){var r=kl(e.domEventName,e.eventSystemFlags,t[0],e.nativeEvent);if(r===null){r=e.nativeEvent;var o=new r.constructor(r.type,r);cl=o,r.target.dispatchEvent(o),cl=null}else return t=po(r),t!==null&&xl(t),e.blockedOn=r,!1;t.shift()}return!0}function hc(e,t,r){Vo(e)&&r.delete(t)}function kp(){wl=!1,br!==null&&Vo(br)&&(br=null),kr!==null&&Vo(kr)&&(kr=null),Sr!==null&&Vo(Sr)&&(Sr=null),qn.forEach(hc),Xn.forEach(hc)}function Jn(e,t){e.blockedOn===t&&(e.blockedOn=null,wl||(wl=!0,l.unstable_scheduleCallback(l.unstable_NormalPriority,kp)))}function eo(e){function t(i){return Jn(i,e)}if(0<Uo.length){Jn(Uo[0],e);for(var r=1;r<Uo.length;r++){var o=Uo[r];o.blockedOn===e&&(o.blockedOn=null)}}for(br!==null&&Jn(br,e),kr!==null&&Jn(kr,e),Sr!==null&&Jn(Sr,e),qn.forEach(t),Xn.forEach(t),r=0;r<Cr.length;r++)o=Cr[r],o.blockedOn===e&&(o.blockedOn=null);for(;0<Cr.length&&(r=Cr[0],r.blockedOn===null);)pc(r),r.blockedOn===null&&Cr.shift()}var hn=J.ReactCurrentBatchConfig,Qo=!0;function Sp(e,t,r,o){var i=Ve,a=hn.transition;hn.transition=null;try{Ve=1,bl(e,t,r,o)}finally{Ve=i,hn.transition=a}}function Cp(e,t,r,o){var i=Ve,a=hn.transition;hn.transition=null;try{Ve=4,bl(e,t,r,o)}finally{Ve=i,hn.transition=a}}function bl(e,t,r,o){if(Qo){var i=kl(e,t,r,o);if(i===null)Nl(e,t,o,Go,r),fc(e,o);else if(bp(i,e,t,r,o))o.stopPropagation();else if(fc(e,o),t&4&&-1<wp.indexOf(e)){for(;i!==null;){var a=po(i);if(a!==null&&sc(a),a=kl(e,t,r,o),a===null&&Nl(e,t,o,Go,r),a===i)break;i=a}i!==null&&o.stopPropagation()}else Nl(e,t,o,null,r)}}var Go=null;function kl(e,t,r,o){if(Go=null,e=ul(o),e=Vr(e),e!==null)if(t=Ur(e),t===null)e=null;else if(r=t.tag,r===13){if(e=Zs(t),e!==null)return e;e=null}else if(r===3){if(t.stateNode.current.memoizedState.isDehydrated)return t.tag===3?t.stateNode.containerInfo:null;e=null}else t!==e&&(e=null);return Go=e,null}function mc(e){switch(e){case"cancel":case"click":case"close":case"contextmenu":case"copy":case"cut":case"auxclick":case"dblclick":case"dragend":case"dragstart":case"drop":case"focusin":case"focusout":case"input":case"invalid":case"keydown":case"keypress":case"keyup":case"mousedown":case"mouseup":case"paste":case"pause":case"play":case"pointercancel":case"pointerdown":case"pointerup":case"ratechange":case"reset":case"resize":case"seeked":case"submit":case"touchcancel":case"touchend":case"touchstart":case"volumechange":case"change":case"selectionchange":case"textInput":case"compositionstart":case"compositionend":case"compositionupdate":case"beforeblur":case"afterblur":case"beforeinput":case"blur":case"fullscreenchange":case"focus":case"hashchange":case"popstate":case"select":case"selectstart":return 1;case"drag":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"mousemove":case"mouseout":case"mouseover":case"pointermove":case"pointerout":case"pointerover":case"scroll":case"toggle":case"touchmove":case"wheel":case"mouseenter":case"mouseleave":case"pointerenter":case"pointerleave":return 4;case"message":switch(dp()){case ml:return 1;case oc:return 4;case No:case fp:return 16;case ic:return 536870912;default:return 16}default:return 16}}var Ar=null,Sl=null,Ko=null;function gc(){if(Ko)return Ko;var e,t=Sl,r=t.length,o,i="value"in Ar?Ar.value:Ar.textContent,a=i.length;for(e=0;e<r&&t[e]===i[e];e++);var d=r-e;for(o=1;o<=d&&t[r-o]===i[a-o];o++);return Ko=i.slice(e,1<o?1-o:void 0)}function Yo(e){var t=e.keyCode;return"charCode"in e?(e=e.charCode,e===0&&t===13&&(e=13)):e=t,e===10&&(e=13),32<=e||e===13?e:0}function qo(){return!0}function vc(){return!1}function $t(e){function t(r,o,i,a,d){this._reactName=r,this._targetInst=i,this.type=o,this.nativeEvent=a,this.target=d,this.currentTarget=null;for(var p in e)e.hasOwnProperty(p)&&(r=e[p],this[p]=r?r(a):a[p]);return this.isDefaultPrevented=(a.defaultPrevented!=null?a.defaultPrevented:a.returnValue===!1)?qo:vc,this.isPropagationStopped=vc,this}return G(t.prototype,{preventDefault:function(){this.defaultPrevented=!0;var r=this.nativeEvent;r&&(r.preventDefault?r.preventDefault():typeof r.returnValue!="unknown"&&(r.returnValue=!1),this.isDefaultPrevented=qo)},stopPropagation:function(){var r=this.nativeEvent;r&&(r.stopPropagation?r.stopPropagation():typeof r.cancelBubble!="unknown"&&(r.cancelBubble=!0),this.isPropagationStopped=qo)},persist:function(){},isPersistent:qo}),t}var mn={eventPhase:0,bubbles:0,cancelable:0,timeStamp:function(e){return e.timeStamp||Date.now()},defaultPrevented:0,isTrusted:0},Cl=$t(mn),to=G({},mn,{view:0,detail:0}),Ap=$t(to),Al,El,ro,Xo=G({},to,{screenX:0,screenY:0,clientX:0,clientY:0,pageX:0,pageY:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,getModifierState:Tl,button:0,buttons:0,relatedTarget:function(e){return e.relatedTarget===void 0?e.fromElement===e.srcElement?e.toElement:e.fromElement:e.relatedTarget},movementX:function(e){return"movementX"in e?e.movementX:(e!==ro&&(ro&&e.type==="mousemove"?(Al=e.screenX-ro.screenX,El=e.screenY-ro.screenY):El=Al=0,ro=e),Al)},movementY:function(e){return"movementY"in e?e.movementY:El}}),yc=$t(Xo),Ep=G({},Xo,{dataTransfer:0}),_p=$t(Ep),Tp=G({},to,{relatedTarget:0}),_l=$t(Tp),jp=G({},mn,{animationName:0,elapsedTime:0,pseudoElement:0}),Pp=$t(jp),Rp=G({},mn,{clipboardData:function(e){return"clipboardData"in e?e.clipboardData:window.clipboardData}}),$p=$t(Rp),Dp=G({},mn,{data:0}),xc=$t(Dp),zp={Esc:"Escape",Spacebar:" ",Left:"ArrowLeft",Up:"ArrowUp",Right:"ArrowRight",Down:"ArrowDown",Del:"Delete",Win:"OS",Menu:"ContextMenu",Apps:"ContextMenu",Scroll:"ScrollLock",MozPrintableKey:"Unidentified"},Ip={8:"Backspace",9:"Tab",12:"Clear",13:"Enter",16:"Shift",17:"Control",18:"Alt",19:"Pause",20:"CapsLock",27:"Escape",32:" ",33:"PageUp",34:"PageDown",35:"End",36:"Home",37:"ArrowLeft",38:"ArrowUp",39:"ArrowRight",40:"ArrowDown",45:"Insert",46:"Delete",112:"F1",113:"F2",114:"F3",115:"F4",116:"F5",117:"F6",118:"F7",119:"F8",120:"F9",121:"F10",122:"F11",123:"F12",144:"NumLock",145:"ScrollLock",224:"Meta"},Lp={Alt:"altKey",Control:"ctrlKey",Meta:"metaKey",Shift:"shiftKey"};function Mp(e){var t=this.nativeEvent;return t.getModifierState?t.getModifierState(e):(e=Lp[e])?!!t[e]:!1}function Tl(){return Mp}var Op=G({},to,{key:function(e){if(e.key){var t=zp[e.key]||e.key;if(t!=="Unidentified")return t}return e.type==="keypress"?(e=Yo(e),e===13?"Enter":String.fromCharCode(e)):e.type==="keydown"||e.type==="keyup"?Ip[e.keyCode]||"Unidentified":""},code:0,location:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,repeat:0,locale:0,getModifierState:Tl,charCode:function(e){return e.type==="keypress"?Yo(e):0},keyCode:function(e){return e.type==="keydown"||e.type==="keyup"?e.keyCode:0},which:function(e){return e.type==="keypress"?Yo(e):e.type==="keydown"||e.type==="keyup"?e.keyCode:0}}),Np=$t(Op),Bp=G({},Xo,{pointerId:0,width:0,height:0,pressure:0,tangentialPressure:0,tiltX:0,tiltY:0,twist:0,pointerType:0,isPrimary:0}),wc=$t(Bp),Fp=G({},to,{touches:0,targetTouches:0,changedTouches:0,altKey:0,metaKey:0,ctrlKey:0,shiftKey:0,getModifierState:Tl}),Wp=$t(Fp),Hp=G({},mn,{propertyName:0,elapsedTime:0,pseudoElement:0}),Up=$t(Hp),Vp=G({},Xo,{deltaX:function(e){return"deltaX"in e?e.deltaX:"wheelDeltaX"in e?-e.wheelDeltaX:0},deltaY:function(e){return"deltaY"in e?e.deltaY:"wheelDeltaY"in e?-e.wheelDeltaY:"wheelDelta"in e?-e.wheelDelta:0},deltaZ:0,deltaMode:0}),Qp=$t(Vp),Gp=[9,13,27,32],jl=v&&"CompositionEvent"in window,no=null;v&&"documentMode"in document&&(no=document.documentMode);var Kp=v&&"TextEvent"in window&&!no,bc=v&&(!jl||no&&8<no&&11>=no),kc=" ",Sc=!1;function Cc(e,t){switch(e){case"keyup":return Gp.indexOf(t.keyCode)!==-1;case"keydown":return t.keyCode!==229;case"keypress":case"mousedown":case"focusout":return!0;default:return!1}}function Ac(e){return e=e.detail,typeof e=="object"&&"data"in e?e.data:null}var gn=!1;function Yp(e,t){switch(e){case"compositionend":return Ac(t);case"keypress":return t.which!==32?null:(Sc=!0,kc);case"textInput":return e=t.data,e===kc&&Sc?null:e;default:return null}}function qp(e,t){if(gn)return e==="compositionend"||!jl&&Cc(e,t)?(e=gc(),Ko=Sl=Ar=null,gn=!1,e):null;switch(e){case"paste":return null;case"keypress":if(!(t.ctrlKey||t.altKey||t.metaKey)||t.ctrlKey&&t.altKey){if(t.char&&1<t.char.length)return t.char;if(t.which)return String.fromCharCode(t.which)}return null;case"compositionend":return bc&&t.locale!=="ko"?null:t.data;default:return null}}var Xp={color:!0,date:!0,datetime:!0,"datetime-local":!0,email:!0,month:!0,number:!0,password:!0,range:!0,search:!0,tel:!0,text:!0,time:!0,url:!0,week:!0};function Ec(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t==="input"?!!Xp[e.type]:t==="textarea"}function _c(e,t,r,o){Gs(o),t=ri(t,"onChange"),0<t.length&&(r=new Cl("onChange","change",null,r,o),e.push({event:r,listeners:t}))}var oo=null,io=null;function Zp(e){Vc(e,0)}function Zo(e){var t=bn(e);if(ae(t))return e}function Jp(e,t){if(e==="change")return t}var Tc=!1;if(v){var Pl;if(v){var Rl="oninput"in document;if(!Rl){var jc=document.createElement("div");jc.setAttribute("oninput","return;"),Rl=typeof jc.oninput=="function"}Pl=Rl}else Pl=!1;Tc=Pl&&(!document.documentMode||9<document.documentMode)}function Pc(){oo&&(oo.detachEvent("onpropertychange",Rc),io=oo=null)}function Rc(e){if(e.propertyName==="value"&&Zo(io)){var t=[];_c(t,io,e,ul(e)),Xs(Zp,t)}}function eh(e,t,r){e==="focusin"?(Pc(),oo=t,io=r,oo.attachEvent("onpropertychange",Rc)):e==="focusout"&&Pc()}function th(e){if(e==="selectionchange"||e==="keyup"||e==="keydown")return Zo(io)}function rh(e,t){if(e==="click")return Zo(t)}function nh(e,t){if(e==="input"||e==="change")return Zo(t)}function oh(e,t){return e===t&&(e!==0||1/e===1/t)||e!==e&&t!==t}var Qt=typeof Object.is=="function"?Object.is:oh;function lo(e,t){if(Qt(e,t))return!0;if(typeof e!="object"||e===null||typeof t!="object"||t===null)return!1;var r=Object.keys(e),o=Object.keys(t);if(r.length!==o.length)return!1;for(o=0;o<r.length;o++){var i=r[o];if(!x.call(t,i)||!Qt(e[i],t[i]))return!1}return!0}function $c(e){for(;e&&e.firstChild;)e=e.firstChild;return e}function Dc(e,t){var r=$c(e);e=0;for(var o;r;){if(r.nodeType===3){if(o=e+r.textContent.length,e<=t&&o>=t)return{node:r,offset:t-e};e=o}e:{for(;r;){if(r.nextSibling){r=r.nextSibling;break e}r=r.parentNode}r=void 0}r=$c(r)}}function zc(e,t){return e&&t?e===t?!0:e&&e.nodeType===3?!1:t&&t.nodeType===3?zc(e,t.parentNode):"contains"in e?e.contains(t):e.compareDocumentPosition?!!(e.compareDocumentPosition(t)&16):!1:!1}function Ic(){for(var e=window,t=je();t instanceof e.HTMLIFrameElement;){try{var r=typeof t.contentWindow.location.href=="string"}catch{r=!1}if(r)e=t.contentWindow;else break;t=je(e.document)}return t}function $l(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t&&(t==="input"&&(e.type==="text"||e.type==="search"||e.type==="tel"||e.type==="url"||e.type==="password")||t==="textarea"||e.contentEditable==="true")}function ih(e){var t=Ic(),r=e.focusedElem,o=e.selectionRange;if(t!==r&&r&&r.ownerDocument&&zc(r.ownerDocument.documentElement,r)){if(o!==null&&$l(r)){if(t=o.start,e=o.end,e===void 0&&(e=t),"selectionStart"in r)r.selectionStart=t,r.selectionEnd=Math.min(e,r.value.length);else if(e=(t=r.ownerDocument||document)&&t.defaultView||window,e.getSelection){e=e.getSelection();var i=r.textContent.length,a=Math.min(o.start,i);o=o.end===void 0?a:Math.min(o.end,i),!e.extend&&a>o&&(i=o,o=a,a=i),i=Dc(r,a);var d=Dc(r,o);i&&d&&(e.rangeCount!==1||e.anchorNode!==i.node||e.anchorOffset!==i.offset||e.focusNode!==d.node||e.focusOffset!==d.offset)&&(t=t.createRange(),t.setStart(i.node,i.offset),e.removeAllRanges(),a>o?(e.addRange(t),e.extend(d.node,d.offset)):(t.setEnd(d.node,d.offset),e.addRange(t)))}}for(t=[],e=r;e=e.parentNode;)e.nodeType===1&&t.push({element:e,left:e.scrollLeft,top:e.scrollTop});for(typeof r.focus=="function"&&r.focus(),r=0;r<t.length;r++)e=t[r],e.element.scrollLeft=e.left,e.element.scrollTop=e.top}}var lh=v&&"documentMode"in document&&11>=document.documentMode,vn=null,Dl=null,ao=null,zl=!1;function Lc(e,t,r){var o=r.window===r?r.document:r.nodeType===9?r:r.ownerDocument;zl||vn==null||vn!==je(o)||(o=vn,"selectionStart"in o&&$l(o)?o={start:o.selectionStart,end:o.selectionEnd}:(o=(o.ownerDocument&&o.ownerDocument.defaultView||window).getSelection(),o={anchorNode:o.anchorNode,anchorOffset:o.anchorOffset,focusNode:o.focusNode,focusOffset:o.focusOffset}),ao&&lo(ao,o)||(ao=o,o=ri(Dl,"onSelect"),0<o.length&&(t=new Cl("onSelect","select",null,t,r),e.push({event:t,listeners:o}),t.target=vn)))}function Jo(e,t){var r={};return r[e.toLowerCase()]=t.toLowerCase(),r["Webkit"+e]="webkit"+t,r["Moz"+e]="moz"+t,r}var yn={animationend:Jo("Animation","AnimationEnd"),animationiteration:Jo("Animation","AnimationIteration"),animationstart:Jo("Animation","AnimationStart"),transitionend:Jo("Transition","TransitionEnd")},Il={},Mc={};v&&(Mc=document.createElement("div").style,"AnimationEvent"in window||(delete yn.animationend.animation,delete yn.animationiteration.animation,delete yn.animationstart.animation),"TransitionEvent"in window||delete yn.transitionend.transition);function ei(e){if(Il[e])return Il[e];if(!yn[e])return e;var t=yn[e],r;for(r in t)if(t.hasOwnProperty(r)&&r in Mc)return Il[e]=t[r];return e}var Oc=ei("animationend"),Nc=ei("animationiteration"),Bc=ei("animationstart"),Fc=ei("transitionend"),Wc=new Map,Hc="abort auxClick cancel canPlay canPlayThrough click close contextMenu copy cut drag dragEnd dragEnter dragExit dragLeave dragOver dragStart drop durationChange emptied encrypted ended error gotPointerCapture input invalid keyDown keyPress keyUp load loadedData loadedMetadata loadStart lostPointerCapture mouseDown mouseMove mouseOut mouseOver mouseUp paste pause play playing pointerCancel pointerDown pointerMove pointerOut pointerOver pointerUp progress rateChange reset resize seeked seeking stalled submit suspend timeUpdate touchCancel touchEnd touchStart volumeChange scroll toggle touchMove waiting wheel".split(" ");function Er(e,t){Wc.set(e,t),g(t,[e])}for(var Ll=0;Ll<Hc.length;Ll++){var Ml=Hc[Ll],ah=Ml.toLowerCase(),sh=Ml[0].toUpperCase()+Ml.slice(1);Er(ah,"on"+sh)}Er(Oc,"onAnimationEnd"),Er(Nc,"onAnimationIteration"),Er(Bc,"onAnimationStart"),Er("dblclick","onDoubleClick"),Er("focusin","onFocus"),Er("focusout","onBlur"),Er(Fc,"onTransitionEnd"),h("onMouseEnter",["mouseout","mouseover"]),h("onMouseLeave",["mouseout","mouseover"]),h("onPointerEnter",["pointerout","pointerover"]),h("onPointerLeave",["pointerout","pointerover"]),g("onChange","change click focusin focusout input keydown keyup selectionchange".split(" ")),g("onSelect","focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange".split(" ")),g("onBeforeInput",["compositionend","keypress","textInput","paste"]),g("onCompositionEnd","compositionend focusout keydown keypress keyup mousedown".split(" ")),g("onCompositionStart","compositionstart focusout keydown keypress keyup mousedown".split(" ")),g("onCompositionUpdate","compositionupdate focusout keydown keypress keyup mousedown".split(" "));var so="abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange resize seeked seeking stalled suspend timeupdate volumechange waiting".split(" "),ch=new Set("cancel close invalid load scroll toggle".split(" ").concat(so));function Uc(e,t,r){var o=e.type||"unknown-event";e.currentTarget=r,ap(o,t,void 0,e),e.currentTarget=null}function Vc(e,t){t=(t&4)!==0;for(var r=0;r<e.length;r++){var o=e[r],i=o.event;o=o.listeners;e:{var a=void 0;if(t)for(var d=o.length-1;0<=d;d--){var p=o[d],w=p.instance,O=p.currentTarget;if(p=p.listener,w!==a&&i.isPropagationStopped())break e;Uc(i,p,O),a=w}else for(d=0;d<o.length;d++){if(p=o[d],w=p.instance,O=p.currentTarget,p=p.listener,w!==a&&i.isPropagationStopped())break e;Uc(i,p,O),a=w}}}if(Oo)throw e=hl,Oo=!1,hl=null,e}function Ye(e,t){var r=t[Vl];r===void 0&&(r=t[Vl]=new Set);var o=e+"__bubble";r.has(o)||(Qc(t,e,2,!1),r.add(o))}function Ol(e,t,r){var o=0;t&&(o|=4),Qc(r,e,o,t)}var ti="_reactListening"+Math.random().toString(36).slice(2);function co(e){if(!e[ti]){e[ti]=!0,u.forEach(function(r){r!=="selectionchange"&&(ch.has(r)||Ol(r,!1,e),Ol(r,!0,e))});var t=e.nodeType===9?e:e.ownerDocument;t===null||t[ti]||(t[ti]=!0,Ol("selectionchange",!1,t))}}function Qc(e,t,r,o){switch(mc(t)){case 1:var i=Sp;break;case 4:i=Cp;break;default:i=bl}r=i.bind(null,t,r,e),i=void 0,!pl||t!=="touchstart"&&t!=="touchmove"&&t!=="wheel"||(i=!0),o?i!==void 0?e.addEventListener(t,r,{capture:!0,passive:i}):e.addEventListener(t,r,!0):i!==void 0?e.addEventListener(t,r,{passive:i}):e.addEventListener(t,r,!1)}function Nl(e,t,r,o,i){var a=o;if((t&1)===0&&(t&2)===0&&o!==null)e:for(;;){if(o===null)return;var d=o.tag;if(d===3||d===4){var p=o.stateNode.containerInfo;if(p===i||p.nodeType===8&&p.parentNode===i)break;if(d===4)for(d=o.return;d!==null;){var w=d.tag;if((w===3||w===4)&&(w=d.stateNode.containerInfo,w===i||w.nodeType===8&&w.parentNode===i))return;d=d.return}for(;p!==null;){if(d=Vr(p),d===null)return;if(w=d.tag,w===5||w===6){o=a=d;continue e}p=p.parentNode}}o=o.return}Xs(function(){var O=a,Z=ul(r),re=[];e:{var X=Wc.get(e);if(X!==void 0){var pe=Cl,we=e;switch(e){case"keypress":if(Yo(r)===0)break e;case"keydown":case"keyup":pe=Np;break;case"focusin":we="focus",pe=_l;break;case"focusout":we="blur",pe=_l;break;case"beforeblur":case"afterblur":pe=_l;break;case"click":if(r.button===2)break e;case"auxclick":case"dblclick":case"mousedown":case"mousemove":case"mouseup":case"mouseout":case"mouseover":case"contextmenu":pe=yc;break;case"drag":case"dragend":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"dragstart":case"drop":pe=_p;break;case"touchcancel":case"touchend":case"touchmove":case"touchstart":pe=Wp;break;case Oc:case Nc:case Bc:pe=Pp;break;case Fc:pe=Up;break;case"scroll":pe=Ap;break;case"wheel":pe=Qp;break;case"copy":case"cut":case"paste":pe=$p;break;case"gotpointercapture":case"lostpointercapture":case"pointercancel":case"pointerdown":case"pointermove":case"pointerout":case"pointerover":case"pointerup":pe=wc}var ke=(t&4)!==0,lt=!ke&&e==="scroll",R=ke?X!==null?X+"Capture":null:X;ke=[];for(var k=O,z;k!==null;){z=k;var ie=z.stateNode;if(z.tag===5&&ie!==null&&(z=ie,R!==null&&(ie=Vn(k,R),ie!=null&&ke.push(uo(k,ie,z)))),lt)break;k=k.return}0<ke.length&&(X=new pe(X,we,null,r,Z),re.push({event:X,listeners:ke}))}}if((t&7)===0){e:{if(X=e==="mouseover"||e==="pointerover",pe=e==="mouseout"||e==="pointerout",X&&r!==cl&&(we=r.relatedTarget||r.fromElement)&&(Vr(we)||we[dr]))break e;if((pe||X)&&(X=Z.window===Z?Z:(X=Z.ownerDocument)?X.defaultView||X.parentWindow:window,pe?(we=r.relatedTarget||r.toElement,pe=O,we=we?Vr(we):null,we!==null&&(lt=Ur(we),we!==lt||we.tag!==5&&we.tag!==6)&&(we=null)):(pe=null,we=O),pe!==we)){if(ke=yc,ie="onMouseLeave",R="onMouseEnter",k="mouse",(e==="pointerout"||e==="pointerover")&&(ke=wc,ie="onPointerLeave",R="onPointerEnter",k="pointer"),lt=pe==null?X:bn(pe),z=we==null?X:bn(we),X=new ke(ie,k+"leave",pe,r,Z),X.target=lt,X.relatedTarget=z,ie=null,Vr(Z)===O&&(ke=new ke(R,k+"enter",we,r,Z),ke.target=z,ke.relatedTarget=lt,ie=ke),lt=ie,pe&&we)t:{for(ke=pe,R=we,k=0,z=ke;z;z=xn(z))k++;for(z=0,ie=R;ie;ie=xn(ie))z++;for(;0<k-z;)ke=xn(ke),k--;for(;0<z-k;)R=xn(R),z--;for(;k--;){if(ke===R||R!==null&&ke===R.alternate)break t;ke=xn(ke),R=xn(R)}ke=null}else ke=null;pe!==null&&Gc(re,X,pe,ke,!1),we!==null&&lt!==null&&Gc(re,lt,we,ke,!0)}}e:{if(X=O?bn(O):window,pe=X.nodeName&&X.nodeName.toLowerCase(),pe==="select"||pe==="input"&&X.type==="file")var Se=Jp;else if(Ec(X))if(Tc)Se=nh;else{Se=th;var _e=eh}else(pe=X.nodeName)&&pe.toLowerCase()==="input"&&(X.type==="checkbox"||X.type==="radio")&&(Se=rh);if(Se&&(Se=Se(e,O))){_c(re,Se,r,Z);break e}_e&&_e(e,X,O),e==="focusout"&&(_e=X._wrapperState)&&_e.controlled&&X.type==="number"&&ot(X,"number",X.value)}switch(_e=O?bn(O):window,e){case"focusin":(Ec(_e)||_e.contentEditable==="true")&&(vn=_e,Dl=O,ao=null);break;case"focusout":ao=Dl=vn=null;break;case"mousedown":zl=!0;break;case"contextmenu":case"mouseup":case"dragend":zl=!1,Lc(re,r,Z);break;case"selectionchange":if(lh)break;case"keydown":case"keyup":Lc(re,r,Z)}var Te;if(jl)e:{switch(e){case"compositionstart":var Pe="onCompositionStart";break e;case"compositionend":Pe="onCompositionEnd";break e;case"compositionupdate":Pe="onCompositionUpdate";break e}Pe=void 0}else gn?Cc(e,r)&&(Pe="onCompositionEnd"):e==="keydown"&&r.keyCode===229&&(Pe="onCompositionStart");Pe&&(bc&&r.locale!=="ko"&&(gn||Pe!=="onCompositionStart"?Pe==="onCompositionEnd"&&gn&&(Te=gc()):(Ar=Z,Sl="value"in Ar?Ar.value:Ar.textContent,gn=!0)),_e=ri(O,Pe),0<_e.length&&(Pe=new xc(Pe,e,null,r,Z),re.push({event:Pe,listeners:_e}),Te?Pe.data=Te:(Te=Ac(r),Te!==null&&(Pe.data=Te)))),(Te=Kp?Yp(e,r):qp(e,r))&&(O=ri(O,"onBeforeInput"),0<O.length&&(Z=new xc("onBeforeInput","beforeinput",null,r,Z),re.push({event:Z,listeners:O}),Z.data=Te))}Vc(re,t)})}function uo(e,t,r){return{instance:e,listener:t,currentTarget:r}}function ri(e,t){for(var r=t+"Capture",o=[];e!==null;){var i=e,a=i.stateNode;i.tag===5&&a!==null&&(i=a,a=Vn(e,r),a!=null&&o.unshift(uo(e,a,i)),a=Vn(e,t),a!=null&&o.push(uo(e,a,i))),e=e.return}return o}function xn(e){if(e===null)return null;do e=e.return;while(e&&e.tag!==5);return e||null}function Gc(e,t,r,o,i){for(var a=t._reactName,d=[];r!==null&&r!==o;){var p=r,w=p.alternate,O=p.stateNode;if(w!==null&&w===o)break;p.tag===5&&O!==null&&(p=O,i?(w=Vn(r,a),w!=null&&d.unshift(uo(r,w,p))):i||(w=Vn(r,a),w!=null&&d.push(uo(r,w,p)))),r=r.return}d.length!==0&&e.push({event:t,listeners:d})}var uh=/\r\n?/g,dh=/\u0000|\uFFFD/g;function Kc(e){return(typeof e=="string"?e:""+e).replace(uh,`
`).replace(dh,"")}function ni(e,t,r){if(t=Kc(t),Kc(e)!==t&&r)throw Error(s(425))}function oi(){}var Bl=null,Fl=null;function Wl(e,t){return e==="textarea"||e==="noscript"||typeof t.children=="string"||typeof t.children=="number"||typeof t.dangerouslySetInnerHTML=="object"&&t.dangerouslySetInnerHTML!==null&&t.dangerouslySetInnerHTML.__html!=null}var Hl=typeof setTimeout=="function"?setTimeout:void 0,fh=typeof clearTimeout=="function"?clearTimeout:void 0,Yc=typeof Promise=="function"?Promise:void 0,ph=typeof queueMicrotask=="function"?queueMicrotask:typeof Yc<"u"?function(e){return Yc.resolve(null).then(e).catch(hh)}:Hl;function hh(e){setTimeout(function(){throw e})}function Ul(e,t){var r=t,o=0;do{var i=r.nextSibling;if(e.removeChild(r),i&&i.nodeType===8)if(r=i.data,r==="/$"){if(o===0){e.removeChild(i),eo(t);return}o--}else r!=="$"&&r!=="$?"&&r!=="$!"||o++;r=i}while(r);eo(t)}function _r(e){for(;e!=null;e=e.nextSibling){var t=e.nodeType;if(t===1||t===3)break;if(t===8){if(t=e.data,t==="$"||t==="$!"||t==="$?")break;if(t==="/$")return null}}return e}function qc(e){e=e.previousSibling;for(var t=0;e;){if(e.nodeType===8){var r=e.data;if(r==="$"||r==="$!"||r==="$?"){if(t===0)return e;t--}else r==="/$"&&t++}e=e.previousSibling}return null}var wn=Math.random().toString(36).slice(2),rr="__reactFiber$"+wn,fo="__reactProps$"+wn,dr="__reactContainer$"+wn,Vl="__reactEvents$"+wn,mh="__reactListeners$"+wn,gh="__reactHandles$"+wn;function Vr(e){var t=e[rr];if(t)return t;for(var r=e.parentNode;r;){if(t=r[dr]||r[rr]){if(r=t.alternate,t.child!==null||r!==null&&r.child!==null)for(e=qc(e);e!==null;){if(r=e[rr])return r;e=qc(e)}return t}e=r,r=e.parentNode}return null}function po(e){return e=e[rr]||e[dr],!e||e.tag!==5&&e.tag!==6&&e.tag!==13&&e.tag!==3?null:e}function bn(e){if(e.tag===5||e.tag===6)return e.stateNode;throw Error(s(33))}function ii(e){return e[fo]||null}var Ql=[],kn=-1;function Tr(e){return{current:e}}function qe(e){0>kn||(e.current=Ql[kn],Ql[kn]=null,kn--)}function Ge(e,t){kn++,Ql[kn]=e.current,e.current=t}var jr={},vt=Tr(jr),At=Tr(!1),Qr=jr;function Sn(e,t){var r=e.type.contextTypes;if(!r)return jr;var o=e.stateNode;if(o&&o.__reactInternalMemoizedUnmaskedChildContext===t)return o.__reactInternalMemoizedMaskedChildContext;var i={},a;for(a in r)i[a]=t[a];return o&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=t,e.__reactInternalMemoizedMaskedChildContext=i),i}function Et(e){return e=e.childContextTypes,e!=null}function li(){qe(At),qe(vt)}function Xc(e,t,r){if(vt.current!==jr)throw Error(s(168));Ge(vt,t),Ge(At,r)}function Zc(e,t,r){var o=e.stateNode;if(t=t.childContextTypes,typeof o.getChildContext!="function")return r;o=o.getChildContext();for(var i in o)if(!(i in t))throw Error(s(108,oe(e)||"Unknown",i));return G({},r,o)}function ai(e){return e=(e=e.stateNode)&&e.__reactInternalMemoizedMergedChildContext||jr,Qr=vt.current,Ge(vt,e),Ge(At,At.current),!0}function Jc(e,t,r){var o=e.stateNode;if(!o)throw Error(s(169));r?(e=Zc(e,t,Qr),o.__reactInternalMemoizedMergedChildContext=e,qe(At),qe(vt),Ge(vt,e)):qe(At),Ge(At,r)}var fr=null,si=!1,Gl=!1;function eu(e){fr===null?fr=[e]:fr.push(e)}function vh(e){si=!0,eu(e)}function Pr(){if(!Gl&&fr!==null){Gl=!0;var e=0,t=Ve;try{var r=fr;for(Ve=1;e<r.length;e++){var o=r[e];do o=o(!0);while(o!==null)}fr=null,si=!1}catch(i){throw fr!==null&&(fr=fr.slice(e+1)),rc(ml,Pr),i}finally{Ve=t,Gl=!1}}return null}var Cn=[],An=0,ci=null,ui=0,Lt=[],Mt=0,Gr=null,pr=1,hr="";function Kr(e,t){Cn[An++]=ui,Cn[An++]=ci,ci=e,ui=t}function tu(e,t,r){Lt[Mt++]=pr,Lt[Mt++]=hr,Lt[Mt++]=Gr,Gr=e;var o=pr;e=hr;var i=32-Vt(o)-1;o&=~(1<<i),r+=1;var a=32-Vt(t)+i;if(30<a){var d=i-i%5;a=(o&(1<<d)-1).toString(32),o>>=d,i-=d,pr=1<<32-Vt(t)+i|r<<i|o,hr=a+e}else pr=1<<a|r<<i|o,hr=e}function Kl(e){e.return!==null&&(Kr(e,1),tu(e,1,0))}function Yl(e){for(;e===ci;)ci=Cn[--An],Cn[An]=null,ui=Cn[--An],Cn[An]=null;for(;e===Gr;)Gr=Lt[--Mt],Lt[Mt]=null,hr=Lt[--Mt],Lt[Mt]=null,pr=Lt[--Mt],Lt[Mt]=null}var Dt=null,zt=null,Xe=!1,Gt=null;function ru(e,t){var r=Ft(5,null,null,0);r.elementType="DELETED",r.stateNode=t,r.return=e,t=e.deletions,t===null?(e.deletions=[r],e.flags|=16):t.push(r)}function nu(e,t){switch(e.tag){case 5:var r=e.type;return t=t.nodeType!==1||r.toLowerCase()!==t.nodeName.toLowerCase()?null:t,t!==null?(e.stateNode=t,Dt=e,zt=_r(t.firstChild),!0):!1;case 6:return t=e.pendingProps===""||t.nodeType!==3?null:t,t!==null?(e.stateNode=t,Dt=e,zt=null,!0):!1;case 13:return t=t.nodeType!==8?null:t,t!==null?(r=Gr!==null?{id:pr,overflow:hr}:null,e.memoizedState={dehydrated:t,treeContext:r,retryLane:1073741824},r=Ft(18,null,null,0),r.stateNode=t,r.return=e,e.child=r,Dt=e,zt=null,!0):!1;default:return!1}}function ql(e){return(e.mode&1)!==0&&(e.flags&128)===0}function Xl(e){if(Xe){var t=zt;if(t){var r=t;if(!nu(e,t)){if(ql(e))throw Error(s(418));t=_r(r.nextSibling);var o=Dt;t&&nu(e,t)?ru(o,r):(e.flags=e.flags&-4097|2,Xe=!1,Dt=e)}}else{if(ql(e))throw Error(s(418));e.flags=e.flags&-4097|2,Xe=!1,Dt=e}}}function ou(e){for(e=e.return;e!==null&&e.tag!==5&&e.tag!==3&&e.tag!==13;)e=e.return;Dt=e}function di(e){if(e!==Dt)return!1;if(!Xe)return ou(e),Xe=!0,!1;var t;if((t=e.tag!==3)&&!(t=e.tag!==5)&&(t=e.type,t=t!=="head"&&t!=="body"&&!Wl(e.type,e.memoizedProps)),t&&(t=zt)){if(ql(e))throw iu(),Error(s(418));for(;t;)ru(e,t),t=_r(t.nextSibling)}if(ou(e),e.tag===13){if(e=e.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(s(317));e:{for(e=e.nextSibling,t=0;e;){if(e.nodeType===8){var r=e.data;if(r==="/$"){if(t===0){zt=_r(e.nextSibling);break e}t--}else r!=="$"&&r!=="$!"&&r!=="$?"||t++}e=e.nextSibling}zt=null}}else zt=Dt?_r(e.stateNode.nextSibling):null;return!0}function iu(){for(var e=zt;e;)e=_r(e.nextSibling)}function En(){zt=Dt=null,Xe=!1}function Zl(e){Gt===null?Gt=[e]:Gt.push(e)}var yh=J.ReactCurrentBatchConfig;function ho(e,t,r){if(e=r.ref,e!==null&&typeof e!="function"&&typeof e!="object"){if(r._owner){if(r=r._owner,r){if(r.tag!==1)throw Error(s(309));var o=r.stateNode}if(!o)throw Error(s(147,e));var i=o,a=""+e;return t!==null&&t.ref!==null&&typeof t.ref=="function"&&t.ref._stringRef===a?t.ref:(t=function(d){var p=i.refs;d===null?delete p[a]:p[a]=d},t._stringRef=a,t)}if(typeof e!="string")throw Error(s(284));if(!r._owner)throw Error(s(290,e))}return e}function fi(e,t){throw e=Object.prototype.toString.call(t),Error(s(31,e==="[object Object]"?"object with keys {"+Object.keys(t).join(", ")+"}":e))}function lu(e){var t=e._init;return t(e._payload)}function au(e){function t(R,k){if(e){var z=R.deletions;z===null?(R.deletions=[k],R.flags|=16):z.push(k)}}function r(R,k){if(!e)return null;for(;k!==null;)t(R,k),k=k.sibling;return null}function o(R,k){for(R=new Map;k!==null;)k.key!==null?R.set(k.key,k):R.set(k.index,k),k=k.sibling;return R}function i(R,k){return R=Or(R,k),R.index=0,R.sibling=null,R}function a(R,k,z){return R.index=z,e?(z=R.alternate,z!==null?(z=z.index,z<k?(R.flags|=2,k):z):(R.flags|=2,k)):(R.flags|=1048576,k)}function d(R){return e&&R.alternate===null&&(R.flags|=2),R}function p(R,k,z,ie){return k===null||k.tag!==6?(k=Ha(z,R.mode,ie),k.return=R,k):(k=i(k,z),k.return=R,k)}function w(R,k,z,ie){var Se=z.type;return Se===ne?Z(R,k,z.props.children,ie,z.key):k!==null&&(k.elementType===Se||typeof Se=="object"&&Se!==null&&Se.$$typeof===ue&&lu(Se)===k.type)?(ie=i(k,z.props),ie.ref=ho(R,k,z),ie.return=R,ie):(ie=Li(z.type,z.key,z.props,null,R.mode,ie),ie.ref=ho(R,k,z),ie.return=R,ie)}function O(R,k,z,ie){return k===null||k.tag!==4||k.stateNode.containerInfo!==z.containerInfo||k.stateNode.implementation!==z.implementation?(k=Ua(z,R.mode,ie),k.return=R,k):(k=i(k,z.children||[]),k.return=R,k)}function Z(R,k,z,ie,Se){return k===null||k.tag!==7?(k=rn(z,R.mode,ie,Se),k.return=R,k):(k=i(k,z),k.return=R,k)}function re(R,k,z){if(typeof k=="string"&&k!==""||typeof k=="number")return k=Ha(""+k,R.mode,z),k.return=R,k;if(typeof k=="object"&&k!==null){switch(k.$$typeof){case L:return z=Li(k.type,k.key,k.props,null,R.mode,z),z.ref=ho(R,null,k),z.return=R,z;case M:return k=Ua(k,R.mode,z),k.return=R,k;case ue:var ie=k._init;return re(R,ie(k._payload),z)}if(Fe(k)||se(k))return k=rn(k,R.mode,z,null),k.return=R,k;fi(R,k)}return null}function X(R,k,z,ie){var Se=k!==null?k.key:null;if(typeof z=="string"&&z!==""||typeof z=="number")return Se!==null?null:p(R,k,""+z,ie);if(typeof z=="object"&&z!==null){switch(z.$$typeof){case L:return z.key===Se?w(R,k,z,ie):null;case M:return z.key===Se?O(R,k,z,ie):null;case ue:return Se=z._init,X(R,k,Se(z._payload),ie)}if(Fe(z)||se(z))return Se!==null?null:Z(R,k,z,ie,null);fi(R,z)}return null}function pe(R,k,z,ie,Se){if(typeof ie=="string"&&ie!==""||typeof ie=="number")return R=R.get(z)||null,p(k,R,""+ie,Se);if(typeof ie=="object"&&ie!==null){switch(ie.$$typeof){case L:return R=R.get(ie.key===null?z:ie.key)||null,w(k,R,ie,Se);case M:return R=R.get(ie.key===null?z:ie.key)||null,O(k,R,ie,Se);case ue:var _e=ie._init;return pe(R,k,z,_e(ie._payload),Se)}if(Fe(ie)||se(ie))return R=R.get(z)||null,Z(k,R,ie,Se,null);fi(k,ie)}return null}function we(R,k,z,ie){for(var Se=null,_e=null,Te=k,Pe=k=0,pt=null;Te!==null&&Pe<z.length;Pe++){Te.index>Pe?(pt=Te,Te=null):pt=Te.sibling;var Be=X(R,Te,z[Pe],ie);if(Be===null){Te===null&&(Te=pt);break}e&&Te&&Be.alternate===null&&t(R,Te),k=a(Be,k,Pe),_e===null?Se=Be:_e.sibling=Be,_e=Be,Te=pt}if(Pe===z.length)return r(R,Te),Xe&&Kr(R,Pe),Se;if(Te===null){for(;Pe<z.length;Pe++)Te=re(R,z[Pe],ie),Te!==null&&(k=a(Te,k,Pe),_e===null?Se=Te:_e.sibling=Te,_e=Te);return Xe&&Kr(R,Pe),Se}for(Te=o(R,Te);Pe<z.length;Pe++)pt=pe(Te,R,Pe,z[Pe],ie),pt!==null&&(e&&pt.alternate!==null&&Te.delete(pt.key===null?Pe:pt.key),k=a(pt,k,Pe),_e===null?Se=pt:_e.sibling=pt,_e=pt);return e&&Te.forEach(function(Nr){return t(R,Nr)}),Xe&&Kr(R,Pe),Se}function ke(R,k,z,ie){var Se=se(z);if(typeof Se!="function")throw Error(s(150));if(z=Se.call(z),z==null)throw Error(s(151));for(var _e=Se=null,Te=k,Pe=k=0,pt=null,Be=z.next();Te!==null&&!Be.done;Pe++,Be=z.next()){Te.index>Pe?(pt=Te,Te=null):pt=Te.sibling;var Nr=X(R,Te,Be.value,ie);if(Nr===null){Te===null&&(Te=pt);break}e&&Te&&Nr.alternate===null&&t(R,Te),k=a(Nr,k,Pe),_e===null?Se=Nr:_e.sibling=Nr,_e=Nr,Te=pt}if(Be.done)return r(R,Te),Xe&&Kr(R,Pe),Se;if(Te===null){for(;!Be.done;Pe++,Be=z.next())Be=re(R,Be.value,ie),Be!==null&&(k=a(Be,k,Pe),_e===null?Se=Be:_e.sibling=Be,_e=Be);return Xe&&Kr(R,Pe),Se}for(Te=o(R,Te);!Be.done;Pe++,Be=z.next())Be=pe(Te,R,Pe,Be.value,ie),Be!==null&&(e&&Be.alternate!==null&&Te.delete(Be.key===null?Pe:Be.key),k=a(Be,k,Pe),_e===null?Se=Be:_e.sibling=Be,_e=Be);return e&&Te.forEach(function(Xh){return t(R,Xh)}),Xe&&Kr(R,Pe),Se}function lt(R,k,z,ie){if(typeof z=="object"&&z!==null&&z.type===ne&&z.key===null&&(z=z.props.children),typeof z=="object"&&z!==null){switch(z.$$typeof){case L:e:{for(var Se=z.key,_e=k;_e!==null;){if(_e.key===Se){if(Se=z.type,Se===ne){if(_e.tag===7){r(R,_e.sibling),k=i(_e,z.props.children),k.return=R,R=k;break e}}else if(_e.elementType===Se||typeof Se=="object"&&Se!==null&&Se.$$typeof===ue&&lu(Se)===_e.type){r(R,_e.sibling),k=i(_e,z.props),k.ref=ho(R,_e,z),k.return=R,R=k;break e}r(R,_e);break}else t(R,_e);_e=_e.sibling}z.type===ne?(k=rn(z.props.children,R.mode,ie,z.key),k.return=R,R=k):(ie=Li(z.type,z.key,z.props,null,R.mode,ie),ie.ref=ho(R,k,z),ie.return=R,R=ie)}return d(R);case M:e:{for(_e=z.key;k!==null;){if(k.key===_e)if(k.tag===4&&k.stateNode.containerInfo===z.containerInfo&&k.stateNode.implementation===z.implementation){r(R,k.sibling),k=i(k,z.children||[]),k.return=R,R=k;break e}else{r(R,k);break}else t(R,k);k=k.sibling}k=Ua(z,R.mode,ie),k.return=R,R=k}return d(R);case ue:return _e=z._init,lt(R,k,_e(z._payload),ie)}if(Fe(z))return we(R,k,z,ie);if(se(z))return ke(R,k,z,ie);fi(R,z)}return typeof z=="string"&&z!==""||typeof z=="number"?(z=""+z,k!==null&&k.tag===6?(r(R,k.sibling),k=i(k,z),k.return=R,R=k):(r(R,k),k=Ha(z,R.mode,ie),k.return=R,R=k),d(R)):r(R,k)}return lt}var _n=au(!0),su=au(!1),pi=Tr(null),hi=null,Tn=null,Jl=null;function ea(){Jl=Tn=hi=null}function ta(e){var t=pi.current;qe(pi),e._currentValue=t}function ra(e,t,r){for(;e!==null;){var o=e.alternate;if((e.childLanes&t)!==t?(e.childLanes|=t,o!==null&&(o.childLanes|=t)):o!==null&&(o.childLanes&t)!==t&&(o.childLanes|=t),e===r)break;e=e.return}}function jn(e,t){hi=e,Jl=Tn=null,e=e.dependencies,e!==null&&e.firstContext!==null&&((e.lanes&t)!==0&&(_t=!0),e.firstContext=null)}function Ot(e){var t=e._currentValue;if(Jl!==e)if(e={context:e,memoizedValue:t,next:null},Tn===null){if(hi===null)throw Error(s(308));Tn=e,hi.dependencies={lanes:0,firstContext:e}}else Tn=Tn.next=e;return t}var Yr=null;function na(e){Yr===null?Yr=[e]:Yr.push(e)}function cu(e,t,r,o){var i=t.interleaved;return i===null?(r.next=r,na(t)):(r.next=i.next,i.next=r),t.interleaved=r,mr(e,o)}function mr(e,t){e.lanes|=t;var r=e.alternate;for(r!==null&&(r.lanes|=t),r=e,e=e.return;e!==null;)e.childLanes|=t,r=e.alternate,r!==null&&(r.childLanes|=t),r=e,e=e.return;return r.tag===3?r.stateNode:null}var Rr=!1;function oa(e){e.updateQueue={baseState:e.memoizedState,firstBaseUpdate:null,lastBaseUpdate:null,shared:{pending:null,interleaved:null,lanes:0},effects:null}}function uu(e,t){e=e.updateQueue,t.updateQueue===e&&(t.updateQueue={baseState:e.baseState,firstBaseUpdate:e.firstBaseUpdate,lastBaseUpdate:e.lastBaseUpdate,shared:e.shared,effects:e.effects})}function gr(e,t){return{eventTime:e,lane:t,tag:0,payload:null,callback:null,next:null}}function $r(e,t,r){var o=e.updateQueue;if(o===null)return null;if(o=o.shared,(Oe&2)!==0){var i=o.pending;return i===null?t.next=t:(t.next=i.next,i.next=t),o.pending=t,mr(e,r)}return i=o.interleaved,i===null?(t.next=t,na(o)):(t.next=i.next,i.next=t),o.interleaved=t,mr(e,r)}function mi(e,t,r){if(t=t.updateQueue,t!==null&&(t=t.shared,(r&4194240)!==0)){var o=t.lanes;o&=e.pendingLanes,r|=o,t.lanes=r,yl(e,r)}}function du(e,t){var r=e.updateQueue,o=e.alternate;if(o!==null&&(o=o.updateQueue,r===o)){var i=null,a=null;if(r=r.firstBaseUpdate,r!==null){do{var d={eventTime:r.eventTime,lane:r.lane,tag:r.tag,payload:r.payload,callback:r.callback,next:null};a===null?i=a=d:a=a.next=d,r=r.next}while(r!==null);a===null?i=a=t:a=a.next=t}else i=a=t;r={baseState:o.baseState,firstBaseUpdate:i,lastBaseUpdate:a,shared:o.shared,effects:o.effects},e.updateQueue=r;return}e=r.lastBaseUpdate,e===null?r.firstBaseUpdate=t:e.next=t,r.lastBaseUpdate=t}function gi(e,t,r,o){var i=e.updateQueue;Rr=!1;var a=i.firstBaseUpdate,d=i.lastBaseUpdate,p=i.shared.pending;if(p!==null){i.shared.pending=null;var w=p,O=w.next;w.next=null,d===null?a=O:d.next=O,d=w;var Z=e.alternate;Z!==null&&(Z=Z.updateQueue,p=Z.lastBaseUpdate,p!==d&&(p===null?Z.firstBaseUpdate=O:p.next=O,Z.lastBaseUpdate=w))}if(a!==null){var re=i.baseState;d=0,Z=O=w=null,p=a;do{var X=p.lane,pe=p.eventTime;if((o&X)===X){Z!==null&&(Z=Z.next={eventTime:pe,lane:0,tag:p.tag,payload:p.payload,callback:p.callback,next:null});e:{var we=e,ke=p;switch(X=t,pe=r,ke.tag){case 1:if(we=ke.payload,typeof we=="function"){re=we.call(pe,re,X);break e}re=we;break e;case 3:we.flags=we.flags&-65537|128;case 0:if(we=ke.payload,X=typeof we=="function"?we.call(pe,re,X):we,X==null)break e;re=G({},re,X);break e;case 2:Rr=!0}}p.callback!==null&&p.lane!==0&&(e.flags|=64,X=i.effects,X===null?i.effects=[p]:X.push(p))}else pe={eventTime:pe,lane:X,tag:p.tag,payload:p.payload,callback:p.callback,next:null},Z===null?(O=Z=pe,w=re):Z=Z.next=pe,d|=X;if(p=p.next,p===null){if(p=i.shared.pending,p===null)break;X=p,p=X.next,X.next=null,i.lastBaseUpdate=X,i.shared.pending=null}}while(!0);if(Z===null&&(w=re),i.baseState=w,i.firstBaseUpdate=O,i.lastBaseUpdate=Z,t=i.shared.interleaved,t!==null){i=t;do d|=i.lane,i=i.next;while(i!==t)}else a===null&&(i.shared.lanes=0);Zr|=d,e.lanes=d,e.memoizedState=re}}function fu(e,t,r){if(e=t.effects,t.effects=null,e!==null)for(t=0;t<e.length;t++){var o=e[t],i=o.callback;if(i!==null){if(o.callback=null,o=r,typeof i!="function")throw Error(s(191,i));i.call(o)}}}var mo={},nr=Tr(mo),go=Tr(mo),vo=Tr(mo);function qr(e){if(e===mo)throw Error(s(174));return e}function ia(e,t){switch(Ge(vo,t),Ge(go,e),Ge(nr,mo),e=t.nodeType,e){case 9:case 11:t=(t=t.documentElement)?t.namespaceURI:Hr(null,"");break;default:e=e===8?t.parentNode:t,t=e.namespaceURI||null,e=e.tagName,t=Hr(t,e)}qe(nr),Ge(nr,t)}function Pn(){qe(nr),qe(go),qe(vo)}function pu(e){qr(vo.current);var t=qr(nr.current),r=Hr(t,e.type);t!==r&&(Ge(go,e),Ge(nr,r))}function la(e){go.current===e&&(qe(nr),qe(go))}var Ze=Tr(0);function vi(e){for(var t=e;t!==null;){if(t.tag===13){var r=t.memoizedState;if(r!==null&&(r=r.dehydrated,r===null||r.data==="$?"||r.data==="$!"))return t}else if(t.tag===19&&t.memoizedProps.revealOrder!==void 0){if((t.flags&128)!==0)return t}else if(t.child!==null){t.child.return=t,t=t.child;continue}if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return null;t=t.return}t.sibling.return=t.return,t=t.sibling}return null}var aa=[];function sa(){for(var e=0;e<aa.length;e++)aa[e]._workInProgressVersionPrimary=null;aa.length=0}var yi=J.ReactCurrentDispatcher,ca=J.ReactCurrentBatchConfig,Xr=0,Je=null,st=null,dt=null,xi=!1,yo=!1,xo=0,xh=0;function yt(){throw Error(s(321))}function ua(e,t){if(t===null)return!1;for(var r=0;r<t.length&&r<e.length;r++)if(!Qt(e[r],t[r]))return!1;return!0}function da(e,t,r,o,i,a){if(Xr=a,Je=t,t.memoizedState=null,t.updateQueue=null,t.lanes=0,yi.current=e===null||e.memoizedState===null?Sh:Ch,e=r(o,i),yo){a=0;do{if(yo=!1,xo=0,25<=a)throw Error(s(301));a+=1,dt=st=null,t.updateQueue=null,yi.current=Ah,e=r(o,i)}while(yo)}if(yi.current=ki,t=st!==null&&st.next!==null,Xr=0,dt=st=Je=null,xi=!1,t)throw Error(s(300));return e}function fa(){var e=xo!==0;return xo=0,e}function or(){var e={memoizedState:null,baseState:null,baseQueue:null,queue:null,next:null};return dt===null?Je.memoizedState=dt=e:dt=dt.next=e,dt}function Nt(){if(st===null){var e=Je.alternate;e=e!==null?e.memoizedState:null}else e=st.next;var t=dt===null?Je.memoizedState:dt.next;if(t!==null)dt=t,st=e;else{if(e===null)throw Error(s(310));st=e,e={memoizedState:st.memoizedState,baseState:st.baseState,baseQueue:st.baseQueue,queue:st.queue,next:null},dt===null?Je.memoizedState=dt=e:dt=dt.next=e}return dt}function wo(e,t){return typeof t=="function"?t(e):t}function pa(e){var t=Nt(),r=t.queue;if(r===null)throw Error(s(311));r.lastRenderedReducer=e;var o=st,i=o.baseQueue,a=r.pending;if(a!==null){if(i!==null){var d=i.next;i.next=a.next,a.next=d}o.baseQueue=i=a,r.pending=null}if(i!==null){a=i.next,o=o.baseState;var p=d=null,w=null,O=a;do{var Z=O.lane;if((Xr&Z)===Z)w!==null&&(w=w.next={lane:0,action:O.action,hasEagerState:O.hasEagerState,eagerState:O.eagerState,next:null}),o=O.hasEagerState?O.eagerState:e(o,O.action);else{var re={lane:Z,action:O.action,hasEagerState:O.hasEagerState,eagerState:O.eagerState,next:null};w===null?(p=w=re,d=o):w=w.next=re,Je.lanes|=Z,Zr|=Z}O=O.next}while(O!==null&&O!==a);w===null?d=o:w.next=p,Qt(o,t.memoizedState)||(_t=!0),t.memoizedState=o,t.baseState=d,t.baseQueue=w,r.lastRenderedState=o}if(e=r.interleaved,e!==null){i=e;do a=i.lane,Je.lanes|=a,Zr|=a,i=i.next;while(i!==e)}else i===null&&(r.lanes=0);return[t.memoizedState,r.dispatch]}function ha(e){var t=Nt(),r=t.queue;if(r===null)throw Error(s(311));r.lastRenderedReducer=e;var o=r.dispatch,i=r.pending,a=t.memoizedState;if(i!==null){r.pending=null;var d=i=i.next;do a=e(a,d.action),d=d.next;while(d!==i);Qt(a,t.memoizedState)||(_t=!0),t.memoizedState=a,t.baseQueue===null&&(t.baseState=a),r.lastRenderedState=a}return[a,o]}function hu(){}function mu(e,t){var r=Je,o=Nt(),i=t(),a=!Qt(o.memoizedState,i);if(a&&(o.memoizedState=i,_t=!0),o=o.queue,ma(yu.bind(null,r,o,e),[e]),o.getSnapshot!==t||a||dt!==null&&dt.memoizedState.tag&1){if(r.flags|=2048,bo(9,vu.bind(null,r,o,i,t),void 0,null),ft===null)throw Error(s(349));(Xr&30)!==0||gu(r,t,i)}return i}function gu(e,t,r){e.flags|=16384,e={getSnapshot:t,value:r},t=Je.updateQueue,t===null?(t={lastEffect:null,stores:null},Je.updateQueue=t,t.stores=[e]):(r=t.stores,r===null?t.stores=[e]:r.push(e))}function vu(e,t,r,o){t.value=r,t.getSnapshot=o,xu(t)&&wu(e)}function yu(e,t,r){return r(function(){xu(t)&&wu(e)})}function xu(e){var t=e.getSnapshot;e=e.value;try{var r=t();return!Qt(e,r)}catch{return!0}}function wu(e){var t=mr(e,1);t!==null&&Xt(t,e,1,-1)}function bu(e){var t=or();return typeof e=="function"&&(e=e()),t.memoizedState=t.baseState=e,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:wo,lastRenderedState:e},t.queue=e,e=e.dispatch=kh.bind(null,Je,e),[t.memoizedState,e]}function bo(e,t,r,o){return e={tag:e,create:t,destroy:r,deps:o,next:null},t=Je.updateQueue,t===null?(t={lastEffect:null,stores:null},Je.updateQueue=t,t.lastEffect=e.next=e):(r=t.lastEffect,r===null?t.lastEffect=e.next=e:(o=r.next,r.next=e,e.next=o,t.lastEffect=e)),e}function ku(){return Nt().memoizedState}function wi(e,t,r,o){var i=or();Je.flags|=e,i.memoizedState=bo(1|t,r,void 0,o===void 0?null:o)}function bi(e,t,r,o){var i=Nt();o=o===void 0?null:o;var a=void 0;if(st!==null){var d=st.memoizedState;if(a=d.destroy,o!==null&&ua(o,d.deps)){i.memoizedState=bo(t,r,a,o);return}}Je.flags|=e,i.memoizedState=bo(1|t,r,a,o)}function Su(e,t){return wi(8390656,8,e,t)}function ma(e,t){return bi(2048,8,e,t)}function Cu(e,t){return bi(4,2,e,t)}function Au(e,t){return bi(4,4,e,t)}function Eu(e,t){if(typeof t=="function")return e=e(),t(e),function(){t(null)};if(t!=null)return e=e(),t.current=e,function(){t.current=null}}function _u(e,t,r){return r=r!=null?r.concat([e]):null,bi(4,4,Eu.bind(null,t,e),r)}function ga(){}function Tu(e,t){var r=Nt();t=t===void 0?null:t;var o=r.memoizedState;return o!==null&&t!==null&&ua(t,o[1])?o[0]:(r.memoizedState=[e,t],e)}function ju(e,t){var r=Nt();t=t===void 0?null:t;var o=r.memoizedState;return o!==null&&t!==null&&ua(t,o[1])?o[0]:(e=e(),r.memoizedState=[e,t],e)}function Pu(e,t,r){return(Xr&21)===0?(e.baseState&&(e.baseState=!1,_t=!0),e.memoizedState=r):(Qt(r,t)||(r=lc(),Je.lanes|=r,Zr|=r,e.baseState=!0),t)}function wh(e,t){var r=Ve;Ve=r!==0&&4>r?r:4,e(!0);var o=ca.transition;ca.transition={};try{e(!1),t()}finally{Ve=r,ca.transition=o}}function Ru(){return Nt().memoizedState}function bh(e,t,r){var o=Lr(e);if(r={lane:o,action:r,hasEagerState:!1,eagerState:null,next:null},$u(e))Du(t,r);else if(r=cu(e,t,r,o),r!==null){var i=kt();Xt(r,e,o,i),zu(r,t,o)}}function kh(e,t,r){var o=Lr(e),i={lane:o,action:r,hasEagerState:!1,eagerState:null,next:null};if($u(e))Du(t,i);else{var a=e.alternate;if(e.lanes===0&&(a===null||a.lanes===0)&&(a=t.lastRenderedReducer,a!==null))try{var d=t.lastRenderedState,p=a(d,r);if(i.hasEagerState=!0,i.eagerState=p,Qt(p,d)){var w=t.interleaved;w===null?(i.next=i,na(t)):(i.next=w.next,w.next=i),t.interleaved=i;return}}catch{}finally{}r=cu(e,t,i,o),r!==null&&(i=kt(),Xt(r,e,o,i),zu(r,t,o))}}function $u(e){var t=e.alternate;return e===Je||t!==null&&t===Je}function Du(e,t){yo=xi=!0;var r=e.pending;r===null?t.next=t:(t.next=r.next,r.next=t),e.pending=t}function zu(e,t,r){if((r&4194240)!==0){var o=t.lanes;o&=e.pendingLanes,r|=o,t.lanes=r,yl(e,r)}}var ki={readContext:Ot,useCallback:yt,useContext:yt,useEffect:yt,useImperativeHandle:yt,useInsertionEffect:yt,useLayoutEffect:yt,useMemo:yt,useReducer:yt,useRef:yt,useState:yt,useDebugValue:yt,useDeferredValue:yt,useTransition:yt,useMutableSource:yt,useSyncExternalStore:yt,useId:yt,unstable_isNewReconciler:!1},Sh={readContext:Ot,useCallback:function(e,t){return or().memoizedState=[e,t===void 0?null:t],e},useContext:Ot,useEffect:Su,useImperativeHandle:function(e,t,r){return r=r!=null?r.concat([e]):null,wi(4194308,4,Eu.bind(null,t,e),r)},useLayoutEffect:function(e,t){return wi(4194308,4,e,t)},useInsertionEffect:function(e,t){return wi(4,2,e,t)},useMemo:function(e,t){var r=or();return t=t===void 0?null:t,e=e(),r.memoizedState=[e,t],e},useReducer:function(e,t,r){var o=or();return t=r!==void 0?r(t):t,o.memoizedState=o.baseState=t,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:e,lastRenderedState:t},o.queue=e,e=e.dispatch=bh.bind(null,Je,e),[o.memoizedState,e]},useRef:function(e){var t=or();return e={current:e},t.memoizedState=e},useState:bu,useDebugValue:ga,useDeferredValue:function(e){return or().memoizedState=e},useTransition:function(){var e=bu(!1),t=e[0];return e=wh.bind(null,e[1]),or().memoizedState=e,[t,e]},useMutableSource:function(){},useSyncExternalStore:function(e,t,r){var o=Je,i=or();if(Xe){if(r===void 0)throw Error(s(407));r=r()}else{if(r=t(),ft===null)throw Error(s(349));(Xr&30)!==0||gu(o,t,r)}i.memoizedState=r;var a={value:r,getSnapshot:t};return i.queue=a,Su(yu.bind(null,o,a,e),[e]),o.flags|=2048,bo(9,vu.bind(null,o,a,r,t),void 0,null),r},useId:function(){var e=or(),t=ft.identifierPrefix;if(Xe){var r=hr,o=pr;r=(o&~(1<<32-Vt(o)-1)).toString(32)+r,t=":"+t+"R"+r,r=xo++,0<r&&(t+="H"+r.toString(32)),t+=":"}else r=xh++,t=":"+t+"r"+r.toString(32)+":";return e.memoizedState=t},unstable_isNewReconciler:!1},Ch={readContext:Ot,useCallback:Tu,useContext:Ot,useEffect:ma,useImperativeHandle:_u,useInsertionEffect:Cu,useLayoutEffect:Au,useMemo:ju,useReducer:pa,useRef:ku,useState:function(){return pa(wo)},useDebugValue:ga,useDeferredValue:function(e){var t=Nt();return Pu(t,st.memoizedState,e)},useTransition:function(){var e=pa(wo)[0],t=Nt().memoizedState;return[e,t]},useMutableSource:hu,useSyncExternalStore:mu,useId:Ru,unstable_isNewReconciler:!1},Ah={readContext:Ot,useCallback:Tu,useContext:Ot,useEffect:ma,useImperativeHandle:_u,useInsertionEffect:Cu,useLayoutEffect:Au,useMemo:ju,useReducer:ha,useRef:ku,useState:function(){return ha(wo)},useDebugValue:ga,useDeferredValue:function(e){var t=Nt();return st===null?t.memoizedState=e:Pu(t,st.memoizedState,e)},useTransition:function(){var e=ha(wo)[0],t=Nt().memoizedState;return[e,t]},useMutableSource:hu,useSyncExternalStore:mu,useId:Ru,unstable_isNewReconciler:!1};function Kt(e,t){if(e&&e.defaultProps){t=G({},t),e=e.defaultProps;for(var r in e)t[r]===void 0&&(t[r]=e[r]);return t}return t}function va(e,t,r,o){t=e.memoizedState,r=r(o,t),r=r==null?t:G({},t,r),e.memoizedState=r,e.lanes===0&&(e.updateQueue.baseState=r)}var Si={isMounted:function(e){return(e=e._reactInternals)?Ur(e)===e:!1},enqueueSetState:function(e,t,r){e=e._reactInternals;var o=kt(),i=Lr(e),a=gr(o,i);a.payload=t,r!=null&&(a.callback=r),t=$r(e,a,i),t!==null&&(Xt(t,e,i,o),mi(t,e,i))},enqueueReplaceState:function(e,t,r){e=e._reactInternals;var o=kt(),i=Lr(e),a=gr(o,i);a.tag=1,a.payload=t,r!=null&&(a.callback=r),t=$r(e,a,i),t!==null&&(Xt(t,e,i,o),mi(t,e,i))},enqueueForceUpdate:function(e,t){e=e._reactInternals;var r=kt(),o=Lr(e),i=gr(r,o);i.tag=2,t!=null&&(i.callback=t),t=$r(e,i,o),t!==null&&(Xt(t,e,o,r),mi(t,e,o))}};function Iu(e,t,r,o,i,a,d){return e=e.stateNode,typeof e.shouldComponentUpdate=="function"?e.shouldComponentUpdate(o,a,d):t.prototype&&t.prototype.isPureReactComponent?!lo(r,o)||!lo(i,a):!0}function Lu(e,t,r){var o=!1,i=jr,a=t.contextType;return typeof a=="object"&&a!==null?a=Ot(a):(i=Et(t)?Qr:vt.current,o=t.contextTypes,a=(o=o!=null)?Sn(e,i):jr),t=new t(r,a),e.memoizedState=t.state!==null&&t.state!==void 0?t.state:null,t.updater=Si,e.stateNode=t,t._reactInternals=e,o&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=i,e.__reactInternalMemoizedMaskedChildContext=a),t}function Mu(e,t,r,o){e=t.state,typeof t.componentWillReceiveProps=="function"&&t.componentWillReceiveProps(r,o),typeof t.UNSAFE_componentWillReceiveProps=="function"&&t.UNSAFE_componentWillReceiveProps(r,o),t.state!==e&&Si.enqueueReplaceState(t,t.state,null)}function ya(e,t,r,o){var i=e.stateNode;i.props=r,i.state=e.memoizedState,i.refs={},oa(e);var a=t.contextType;typeof a=="object"&&a!==null?i.context=Ot(a):(a=Et(t)?Qr:vt.current,i.context=Sn(e,a)),i.state=e.memoizedState,a=t.getDerivedStateFromProps,typeof a=="function"&&(va(e,t,a,r),i.state=e.memoizedState),typeof t.getDerivedStateFromProps=="function"||typeof i.getSnapshotBeforeUpdate=="function"||typeof i.UNSAFE_componentWillMount!="function"&&typeof i.componentWillMount!="function"||(t=i.state,typeof i.componentWillMount=="function"&&i.componentWillMount(),typeof i.UNSAFE_componentWillMount=="function"&&i.UNSAFE_componentWillMount(),t!==i.state&&Si.enqueueReplaceState(i,i.state,null),gi(e,r,i,o),i.state=e.memoizedState),typeof i.componentDidMount=="function"&&(e.flags|=4194308)}function Rn(e,t){try{var r="",o=t;do r+=q(o),o=o.return;while(o);var i=r}catch(a){i=`
Error generating stack: `+a.message+`
`+a.stack}return{value:e,source:t,stack:i,digest:null}}function xa(e,t,r){return{value:e,source:null,stack:r??null,digest:t??null}}function wa(e,t){try{console.error(t.value)}catch(r){setTimeout(function(){throw r})}}var Eh=typeof WeakMap=="function"?WeakMap:Map;function Ou(e,t,r){r=gr(-1,r),r.tag=3,r.payload={element:null};var o=t.value;return r.callback=function(){Pi||(Pi=!0,Ia=o),wa(e,t)},r}function Nu(e,t,r){r=gr(-1,r),r.tag=3;var o=e.type.getDerivedStateFromError;if(typeof o=="function"){var i=t.value;r.payload=function(){return o(i)},r.callback=function(){wa(e,t)}}var a=e.stateNode;return a!==null&&typeof a.componentDidCatch=="function"&&(r.callback=function(){wa(e,t),typeof o!="function"&&(zr===null?zr=new Set([this]):zr.add(this));var d=t.stack;this.componentDidCatch(t.value,{componentStack:d!==null?d:""})}),r}function Bu(e,t,r){var o=e.pingCache;if(o===null){o=e.pingCache=new Eh;var i=new Set;o.set(t,i)}else i=o.get(t),i===void 0&&(i=new Set,o.set(t,i));i.has(r)||(i.add(r),e=Bh.bind(null,e,t,r),t.then(e,e))}function Fu(e){do{var t;if((t=e.tag===13)&&(t=e.memoizedState,t=t!==null?t.dehydrated!==null:!0),t)return e;e=e.return}while(e!==null);return null}function Wu(e,t,r,o,i){return(e.mode&1)===0?(e===t?e.flags|=65536:(e.flags|=128,r.flags|=131072,r.flags&=-52805,r.tag===1&&(r.alternate===null?r.tag=17:(t=gr(-1,1),t.tag=2,$r(r,t,1))),r.lanes|=1),e):(e.flags|=65536,e.lanes=i,e)}var _h=J.ReactCurrentOwner,_t=!1;function bt(e,t,r,o){t.child=e===null?su(t,null,r,o):_n(t,e.child,r,o)}function Hu(e,t,r,o,i){r=r.render;var a=t.ref;return jn(t,i),o=da(e,t,r,o,a,i),r=fa(),e!==null&&!_t?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~i,vr(e,t,i)):(Xe&&r&&Kl(t),t.flags|=1,bt(e,t,o,i),t.child)}function Uu(e,t,r,o,i){if(e===null){var a=r.type;return typeof a=="function"&&!Wa(a)&&a.defaultProps===void 0&&r.compare===null&&r.defaultProps===void 0?(t.tag=15,t.type=a,Vu(e,t,a,o,i)):(e=Li(r.type,null,o,t,t.mode,i),e.ref=t.ref,e.return=t,t.child=e)}if(a=e.child,(e.lanes&i)===0){var d=a.memoizedProps;if(r=r.compare,r=r!==null?r:lo,r(d,o)&&e.ref===t.ref)return vr(e,t,i)}return t.flags|=1,e=Or(a,o),e.ref=t.ref,e.return=t,t.child=e}function Vu(e,t,r,o,i){if(e!==null){var a=e.memoizedProps;if(lo(a,o)&&e.ref===t.ref)if(_t=!1,t.pendingProps=o=a,(e.lanes&i)!==0)(e.flags&131072)!==0&&(_t=!0);else return t.lanes=e.lanes,vr(e,t,i)}return ba(e,t,r,o,i)}function Qu(e,t,r){var o=t.pendingProps,i=o.children,a=e!==null?e.memoizedState:null;if(o.mode==="hidden")if((t.mode&1)===0)t.memoizedState={baseLanes:0,cachePool:null,transitions:null},Ge(Dn,It),It|=r;else{if((r&1073741824)===0)return e=a!==null?a.baseLanes|r:r,t.lanes=t.childLanes=1073741824,t.memoizedState={baseLanes:e,cachePool:null,transitions:null},t.updateQueue=null,Ge(Dn,It),It|=e,null;t.memoizedState={baseLanes:0,cachePool:null,transitions:null},o=a!==null?a.baseLanes:r,Ge(Dn,It),It|=o}else a!==null?(o=a.baseLanes|r,t.memoizedState=null):o=r,Ge(Dn,It),It|=o;return bt(e,t,i,r),t.child}function Gu(e,t){var r=t.ref;(e===null&&r!==null||e!==null&&e.ref!==r)&&(t.flags|=512,t.flags|=2097152)}function ba(e,t,r,o,i){var a=Et(r)?Qr:vt.current;return a=Sn(t,a),jn(t,i),r=da(e,t,r,o,a,i),o=fa(),e!==null&&!_t?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~i,vr(e,t,i)):(Xe&&o&&Kl(t),t.flags|=1,bt(e,t,r,i),t.child)}function Ku(e,t,r,o,i){if(Et(r)){var a=!0;ai(t)}else a=!1;if(jn(t,i),t.stateNode===null)Ai(e,t),Lu(t,r,o),ya(t,r,o,i),o=!0;else if(e===null){var d=t.stateNode,p=t.memoizedProps;d.props=p;var w=d.context,O=r.contextType;typeof O=="object"&&O!==null?O=Ot(O):(O=Et(r)?Qr:vt.current,O=Sn(t,O));var Z=r.getDerivedStateFromProps,re=typeof Z=="function"||typeof d.getSnapshotBeforeUpdate=="function";re||typeof d.UNSAFE_componentWillReceiveProps!="function"&&typeof d.componentWillReceiveProps!="function"||(p!==o||w!==O)&&Mu(t,d,o,O),Rr=!1;var X=t.memoizedState;d.state=X,gi(t,o,d,i),w=t.memoizedState,p!==o||X!==w||At.current||Rr?(typeof Z=="function"&&(va(t,r,Z,o),w=t.memoizedState),(p=Rr||Iu(t,r,p,o,X,w,O))?(re||typeof d.UNSAFE_componentWillMount!="function"&&typeof d.componentWillMount!="function"||(typeof d.componentWillMount=="function"&&d.componentWillMount(),typeof d.UNSAFE_componentWillMount=="function"&&d.UNSAFE_componentWillMount()),typeof d.componentDidMount=="function"&&(t.flags|=4194308)):(typeof d.componentDidMount=="function"&&(t.flags|=4194308),t.memoizedProps=o,t.memoizedState=w),d.props=o,d.state=w,d.context=O,o=p):(typeof d.componentDidMount=="function"&&(t.flags|=4194308),o=!1)}else{d=t.stateNode,uu(e,t),p=t.memoizedProps,O=t.type===t.elementType?p:Kt(t.type,p),d.props=O,re=t.pendingProps,X=d.context,w=r.contextType,typeof w=="object"&&w!==null?w=Ot(w):(w=Et(r)?Qr:vt.current,w=Sn(t,w));var pe=r.getDerivedStateFromProps;(Z=typeof pe=="function"||typeof d.getSnapshotBeforeUpdate=="function")||typeof d.UNSAFE_componentWillReceiveProps!="function"&&typeof d.componentWillReceiveProps!="function"||(p!==re||X!==w)&&Mu(t,d,o,w),Rr=!1,X=t.memoizedState,d.state=X,gi(t,o,d,i);var we=t.memoizedState;p!==re||X!==we||At.current||Rr?(typeof pe=="function"&&(va(t,r,pe,o),we=t.memoizedState),(O=Rr||Iu(t,r,O,o,X,we,w)||!1)?(Z||typeof d.UNSAFE_componentWillUpdate!="function"&&typeof d.componentWillUpdate!="function"||(typeof d.componentWillUpdate=="function"&&d.componentWillUpdate(o,we,w),typeof d.UNSAFE_componentWillUpdate=="function"&&d.UNSAFE_componentWillUpdate(o,we,w)),typeof d.componentDidUpdate=="function"&&(t.flags|=4),typeof d.getSnapshotBeforeUpdate=="function"&&(t.flags|=1024)):(typeof d.componentDidUpdate!="function"||p===e.memoizedProps&&X===e.memoizedState||(t.flags|=4),typeof d.getSnapshotBeforeUpdate!="function"||p===e.memoizedProps&&X===e.memoizedState||(t.flags|=1024),t.memoizedProps=o,t.memoizedState=we),d.props=o,d.state=we,d.context=w,o=O):(typeof d.componentDidUpdate!="function"||p===e.memoizedProps&&X===e.memoizedState||(t.flags|=4),typeof d.getSnapshotBeforeUpdate!="function"||p===e.memoizedProps&&X===e.memoizedState||(t.flags|=1024),o=!1)}return ka(e,t,r,o,a,i)}function ka(e,t,r,o,i,a){Gu(e,t);var d=(t.flags&128)!==0;if(!o&&!d)return i&&Jc(t,r,!1),vr(e,t,a);o=t.stateNode,_h.current=t;var p=d&&typeof r.getDerivedStateFromError!="function"?null:o.render();return t.flags|=1,e!==null&&d?(t.child=_n(t,e.child,null,a),t.child=_n(t,null,p,a)):bt(e,t,p,a),t.memoizedState=o.state,i&&Jc(t,r,!0),t.child}function Yu(e){var t=e.stateNode;t.pendingContext?Xc(e,t.pendingContext,t.pendingContext!==t.context):t.context&&Xc(e,t.context,!1),ia(e,t.containerInfo)}function qu(e,t,r,o,i){return En(),Zl(i),t.flags|=256,bt(e,t,r,o),t.child}var Sa={dehydrated:null,treeContext:null,retryLane:0};function Ca(e){return{baseLanes:e,cachePool:null,transitions:null}}function Xu(e,t,r){var o=t.pendingProps,i=Ze.current,a=!1,d=(t.flags&128)!==0,p;if((p=d)||(p=e!==null&&e.memoizedState===null?!1:(i&2)!==0),p?(a=!0,t.flags&=-129):(e===null||e.memoizedState!==null)&&(i|=1),Ge(Ze,i&1),e===null)return Xl(t),e=t.memoizedState,e!==null&&(e=e.dehydrated,e!==null)?((t.mode&1)===0?t.lanes=1:e.data==="$!"?t.lanes=8:t.lanes=1073741824,null):(d=o.children,e=o.fallback,a?(o=t.mode,a=t.child,d={mode:"hidden",children:d},(o&1)===0&&a!==null?(a.childLanes=0,a.pendingProps=d):a=Mi(d,o,0,null),e=rn(e,o,r,null),a.return=t,e.return=t,a.sibling=e,t.child=a,t.child.memoizedState=Ca(r),t.memoizedState=Sa,e):Aa(t,d));if(i=e.memoizedState,i!==null&&(p=i.dehydrated,p!==null))return Th(e,t,d,o,p,i,r);if(a){a=o.fallback,d=t.mode,i=e.child,p=i.sibling;var w={mode:"hidden",children:o.children};return(d&1)===0&&t.child!==i?(o=t.child,o.childLanes=0,o.pendingProps=w,t.deletions=null):(o=Or(i,w),o.subtreeFlags=i.subtreeFlags&14680064),p!==null?a=Or(p,a):(a=rn(a,d,r,null),a.flags|=2),a.return=t,o.return=t,o.sibling=a,t.child=o,o=a,a=t.child,d=e.child.memoizedState,d=d===null?Ca(r):{baseLanes:d.baseLanes|r,cachePool:null,transitions:d.transitions},a.memoizedState=d,a.childLanes=e.childLanes&~r,t.memoizedState=Sa,o}return a=e.child,e=a.sibling,o=Or(a,{mode:"visible",children:o.children}),(t.mode&1)===0&&(o.lanes=r),o.return=t,o.sibling=null,e!==null&&(r=t.deletions,r===null?(t.deletions=[e],t.flags|=16):r.push(e)),t.child=o,t.memoizedState=null,o}function Aa(e,t){return t=Mi({mode:"visible",children:t},e.mode,0,null),t.return=e,e.child=t}function Ci(e,t,r,o){return o!==null&&Zl(o),_n(t,e.child,null,r),e=Aa(t,t.pendingProps.children),e.flags|=2,t.memoizedState=null,e}function Th(e,t,r,o,i,a,d){if(r)return t.flags&256?(t.flags&=-257,o=xa(Error(s(422))),Ci(e,t,d,o)):t.memoizedState!==null?(t.child=e.child,t.flags|=128,null):(a=o.fallback,i=t.mode,o=Mi({mode:"visible",children:o.children},i,0,null),a=rn(a,i,d,null),a.flags|=2,o.return=t,a.return=t,o.sibling=a,t.child=o,(t.mode&1)!==0&&_n(t,e.child,null,d),t.child.memoizedState=Ca(d),t.memoizedState=Sa,a);if((t.mode&1)===0)return Ci(e,t,d,null);if(i.data==="$!"){if(o=i.nextSibling&&i.nextSibling.dataset,o)var p=o.dgst;return o=p,a=Error(s(419)),o=xa(a,o,void 0),Ci(e,t,d,o)}if(p=(d&e.childLanes)!==0,_t||p){if(o=ft,o!==null){switch(d&-d){case 4:i=2;break;case 16:i=8;break;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:i=32;break;case 536870912:i=268435456;break;default:i=0}i=(i&(o.suspendedLanes|d))!==0?0:i,i!==0&&i!==a.retryLane&&(a.retryLane=i,mr(e,i),Xt(o,e,i,-1))}return Fa(),o=xa(Error(s(421))),Ci(e,t,d,o)}return i.data==="$?"?(t.flags|=128,t.child=e.child,t=Fh.bind(null,e),i._reactRetry=t,null):(e=a.treeContext,zt=_r(i.nextSibling),Dt=t,Xe=!0,Gt=null,e!==null&&(Lt[Mt++]=pr,Lt[Mt++]=hr,Lt[Mt++]=Gr,pr=e.id,hr=e.overflow,Gr=t),t=Aa(t,o.children),t.flags|=4096,t)}function Zu(e,t,r){e.lanes|=t;var o=e.alternate;o!==null&&(o.lanes|=t),ra(e.return,t,r)}function Ea(e,t,r,o,i){var a=e.memoizedState;a===null?e.memoizedState={isBackwards:t,rendering:null,renderingStartTime:0,last:o,tail:r,tailMode:i}:(a.isBackwards=t,a.rendering=null,a.renderingStartTime=0,a.last=o,a.tail=r,a.tailMode=i)}function Ju(e,t,r){var o=t.pendingProps,i=o.revealOrder,a=o.tail;if(bt(e,t,o.children,r),o=Ze.current,(o&2)!==0)o=o&1|2,t.flags|=128;else{if(e!==null&&(e.flags&128)!==0)e:for(e=t.child;e!==null;){if(e.tag===13)e.memoizedState!==null&&Zu(e,r,t);else if(e.tag===19)Zu(e,r,t);else if(e.child!==null){e.child.return=e,e=e.child;continue}if(e===t)break e;for(;e.sibling===null;){if(e.return===null||e.return===t)break e;e=e.return}e.sibling.return=e.return,e=e.sibling}o&=1}if(Ge(Ze,o),(t.mode&1)===0)t.memoizedState=null;else switch(i){case"forwards":for(r=t.child,i=null;r!==null;)e=r.alternate,e!==null&&vi(e)===null&&(i=r),r=r.sibling;r=i,r===null?(i=t.child,t.child=null):(i=r.sibling,r.sibling=null),Ea(t,!1,i,r,a);break;case"backwards":for(r=null,i=t.child,t.child=null;i!==null;){if(e=i.alternate,e!==null&&vi(e)===null){t.child=i;break}e=i.sibling,i.sibling=r,r=i,i=e}Ea(t,!0,r,null,a);break;case"together":Ea(t,!1,null,null,void 0);break;default:t.memoizedState=null}return t.child}function Ai(e,t){(t.mode&1)===0&&e!==null&&(e.alternate=null,t.alternate=null,t.flags|=2)}function vr(e,t,r){if(e!==null&&(t.dependencies=e.dependencies),Zr|=t.lanes,(r&t.childLanes)===0)return null;if(e!==null&&t.child!==e.child)throw Error(s(153));if(t.child!==null){for(e=t.child,r=Or(e,e.pendingProps),t.child=r,r.return=t;e.sibling!==null;)e=e.sibling,r=r.sibling=Or(e,e.pendingProps),r.return=t;r.sibling=null}return t.child}function jh(e,t,r){switch(t.tag){case 3:Yu(t),En();break;case 5:pu(t);break;case 1:Et(t.type)&&ai(t);break;case 4:ia(t,t.stateNode.containerInfo);break;case 10:var o=t.type._context,i=t.memoizedProps.value;Ge(pi,o._currentValue),o._currentValue=i;break;case 13:if(o=t.memoizedState,o!==null)return o.dehydrated!==null?(Ge(Ze,Ze.current&1),t.flags|=128,null):(r&t.child.childLanes)!==0?Xu(e,t,r):(Ge(Ze,Ze.current&1),e=vr(e,t,r),e!==null?e.sibling:null);Ge(Ze,Ze.current&1);break;case 19:if(o=(r&t.childLanes)!==0,(e.flags&128)!==0){if(o)return Ju(e,t,r);t.flags|=128}if(i=t.memoizedState,i!==null&&(i.rendering=null,i.tail=null,i.lastEffect=null),Ge(Ze,Ze.current),o)break;return null;case 22:case 23:return t.lanes=0,Qu(e,t,r)}return vr(e,t,r)}var ed,_a,td,rd;ed=function(e,t){for(var r=t.child;r!==null;){if(r.tag===5||r.tag===6)e.appendChild(r.stateNode);else if(r.tag!==4&&r.child!==null){r.child.return=r,r=r.child;continue}if(r===t)break;for(;r.sibling===null;){if(r.return===null||r.return===t)return;r=r.return}r.sibling.return=r.return,r=r.sibling}},_a=function(){},td=function(e,t,r,o){var i=e.memoizedProps;if(i!==o){e=t.stateNode,qr(nr.current);var a=null;switch(r){case"input":i=ze(e,i),o=ze(e,o),a=[];break;case"select":i=G({},i,{value:void 0}),o=G({},o,{value:void 0}),a=[];break;case"textarea":i=te(e,i),o=te(e,o),a=[];break;default:typeof i.onClick!="function"&&typeof o.onClick=="function"&&(e.onclick=oi)}al(r,o);var d;r=null;for(O in i)if(!o.hasOwnProperty(O)&&i.hasOwnProperty(O)&&i[O]!=null)if(O==="style"){var p=i[O];for(d in p)p.hasOwnProperty(d)&&(r||(r={}),r[d]="")}else O!=="dangerouslySetInnerHTML"&&O!=="children"&&O!=="suppressContentEditableWarning"&&O!=="suppressHydrationWarning"&&O!=="autoFocus"&&(f.hasOwnProperty(O)?a||(a=[]):(a=a||[]).push(O,null));for(O in o){var w=o[O];if(p=i!=null?i[O]:void 0,o.hasOwnProperty(O)&&w!==p&&(w!=null||p!=null))if(O==="style")if(p){for(d in p)!p.hasOwnProperty(d)||w&&w.hasOwnProperty(d)||(r||(r={}),r[d]="");for(d in w)w.hasOwnProperty(d)&&p[d]!==w[d]&&(r||(r={}),r[d]=w[d])}else r||(a||(a=[]),a.push(O,r)),r=w;else O==="dangerouslySetInnerHTML"?(w=w?w.__html:void 0,p=p?p.__html:void 0,w!=null&&p!==w&&(a=a||[]).push(O,w)):O==="children"?typeof w!="string"&&typeof w!="number"||(a=a||[]).push(O,""+w):O!=="suppressContentEditableWarning"&&O!=="suppressHydrationWarning"&&(f.hasOwnProperty(O)?(w!=null&&O==="onScroll"&&Ye("scroll",e),a||p===w||(a=[])):(a=a||[]).push(O,w))}r&&(a=a||[]).push("style",r);var O=a;(t.updateQueue=O)&&(t.flags|=4)}},rd=function(e,t,r,o){r!==o&&(t.flags|=4)};function ko(e,t){if(!Xe)switch(e.tailMode){case"hidden":t=e.tail;for(var r=null;t!==null;)t.alternate!==null&&(r=t),t=t.sibling;r===null?e.tail=null:r.sibling=null;break;case"collapsed":r=e.tail;for(var o=null;r!==null;)r.alternate!==null&&(o=r),r=r.sibling;o===null?t||e.tail===null?e.tail=null:e.tail.sibling=null:o.sibling=null}}function xt(e){var t=e.alternate!==null&&e.alternate.child===e.child,r=0,o=0;if(t)for(var i=e.child;i!==null;)r|=i.lanes|i.childLanes,o|=i.subtreeFlags&14680064,o|=i.flags&14680064,i.return=e,i=i.sibling;else for(i=e.child;i!==null;)r|=i.lanes|i.childLanes,o|=i.subtreeFlags,o|=i.flags,i.return=e,i=i.sibling;return e.subtreeFlags|=o,e.childLanes=r,t}function Ph(e,t,r){var o=t.pendingProps;switch(Yl(t),t.tag){case 2:case 16:case 15:case 0:case 11:case 7:case 8:case 12:case 9:case 14:return xt(t),null;case 1:return Et(t.type)&&li(),xt(t),null;case 3:return o=t.stateNode,Pn(),qe(At),qe(vt),sa(),o.pendingContext&&(o.context=o.pendingContext,o.pendingContext=null),(e===null||e.child===null)&&(di(t)?t.flags|=4:e===null||e.memoizedState.isDehydrated&&(t.flags&256)===0||(t.flags|=1024,Gt!==null&&(Oa(Gt),Gt=null))),_a(e,t),xt(t),null;case 5:la(t);var i=qr(vo.current);if(r=t.type,e!==null&&t.stateNode!=null)td(e,t,r,o,i),e.ref!==t.ref&&(t.flags|=512,t.flags|=2097152);else{if(!o){if(t.stateNode===null)throw Error(s(166));return xt(t),null}if(e=qr(nr.current),di(t)){o=t.stateNode,r=t.type;var a=t.memoizedProps;switch(o[rr]=t,o[fo]=a,e=(t.mode&1)!==0,r){case"dialog":Ye("cancel",o),Ye("close",o);break;case"iframe":case"object":case"embed":Ye("load",o);break;case"video":case"audio":for(i=0;i<so.length;i++)Ye(so[i],o);break;case"source":Ye("error",o);break;case"img":case"image":case"link":Ye("error",o),Ye("load",o);break;case"details":Ye("toggle",o);break;case"input":Ce(o,a),Ye("invalid",o);break;case"select":o._wrapperState={wasMultiple:!!a.multiple},Ye("invalid",o);break;case"textarea":Ee(o,a),Ye("invalid",o)}al(r,a),i=null;for(var d in a)if(a.hasOwnProperty(d)){var p=a[d];d==="children"?typeof p=="string"?o.textContent!==p&&(a.suppressHydrationWarning!==!0&&ni(o.textContent,p,e),i=["children",p]):typeof p=="number"&&o.textContent!==""+p&&(a.suppressHydrationWarning!==!0&&ni(o.textContent,p,e),i=["children",""+p]):f.hasOwnProperty(d)&&p!=null&&d==="onScroll"&&Ye("scroll",o)}switch(r){case"input":fe(o),Ne(o,a,!0);break;case"textarea":fe(o),Ue(o);break;case"select":case"option":break;default:typeof a.onClick=="function"&&(o.onclick=oi)}o=i,t.updateQueue=o,o!==null&&(t.flags|=4)}else{d=i.nodeType===9?i:i.ownerDocument,e==="http://www.w3.org/1999/xhtml"&&(e=ur(r)),e==="http://www.w3.org/1999/xhtml"?r==="script"?(e=d.createElement("div"),e.innerHTML="<script><\/script>",e=e.removeChild(e.firstChild)):typeof o.is=="string"?e=d.createElement(r,{is:o.is}):(e=d.createElement(r),r==="select"&&(d=e,o.multiple?d.multiple=!0:o.size&&(d.size=o.size))):e=d.createElementNS(e,r),e[rr]=t,e[fo]=o,ed(e,t,!1,!1),t.stateNode=e;e:{switch(d=sl(r,o),r){case"dialog":Ye("cancel",e),Ye("close",e),i=o;break;case"iframe":case"object":case"embed":Ye("load",e),i=o;break;case"video":case"audio":for(i=0;i<so.length;i++)Ye(so[i],e);i=o;break;case"source":Ye("error",e),i=o;break;case"img":case"image":case"link":Ye("error",e),Ye("load",e),i=o;break;case"details":Ye("toggle",e),i=o;break;case"input":Ce(e,o),i=ze(e,o),Ye("invalid",e);break;case"option":i=o;break;case"select":e._wrapperState={wasMultiple:!!o.multiple},i=G({},o,{value:void 0}),Ye("invalid",e);break;case"textarea":Ee(e,o),i=te(e,o),Ye("invalid",e);break;default:i=o}al(r,i),p=i;for(a in p)if(p.hasOwnProperty(a)){var w=p[a];a==="style"?Vs(e,w):a==="dangerouslySetInnerHTML"?(w=w?w.__html:void 0,w!=null&&Hs(e,w)):a==="children"?typeof w=="string"?(r!=="textarea"||w!=="")&&Hn(e,w):typeof w=="number"&&Hn(e,""+w):a!=="suppressContentEditableWarning"&&a!=="suppressHydrationWarning"&&a!=="autoFocus"&&(f.hasOwnProperty(a)?w!=null&&a==="onScroll"&&Ye("scroll",e):w!=null&&V(e,a,w,d))}switch(r){case"input":fe(e),Ne(e,o,!1);break;case"textarea":fe(e),Ue(e);break;case"option":o.value!=null&&e.setAttribute("value",""+Q(o.value));break;case"select":e.multiple=!!o.multiple,a=o.value,a!=null?N(e,!!o.multiple,a,!1):o.defaultValue!=null&&N(e,!!o.multiple,o.defaultValue,!0);break;default:typeof i.onClick=="function"&&(e.onclick=oi)}switch(r){case"button":case"input":case"select":case"textarea":o=!!o.autoFocus;break e;case"img":o=!0;break e;default:o=!1}}o&&(t.flags|=4)}t.ref!==null&&(t.flags|=512,t.flags|=2097152)}return xt(t),null;case 6:if(e&&t.stateNode!=null)rd(e,t,e.memoizedProps,o);else{if(typeof o!="string"&&t.stateNode===null)throw Error(s(166));if(r=qr(vo.current),qr(nr.current),di(t)){if(o=t.stateNode,r=t.memoizedProps,o[rr]=t,(a=o.nodeValue!==r)&&(e=Dt,e!==null))switch(e.tag){case 3:ni(o.nodeValue,r,(e.mode&1)!==0);break;case 5:e.memoizedProps.suppressHydrationWarning!==!0&&ni(o.nodeValue,r,(e.mode&1)!==0)}a&&(t.flags|=4)}else o=(r.nodeType===9?r:r.ownerDocument).createTextNode(o),o[rr]=t,t.stateNode=o}return xt(t),null;case 13:if(qe(Ze),o=t.memoizedState,e===null||e.memoizedState!==null&&e.memoizedState.dehydrated!==null){if(Xe&&zt!==null&&(t.mode&1)!==0&&(t.flags&128)===0)iu(),En(),t.flags|=98560,a=!1;else if(a=di(t),o!==null&&o.dehydrated!==null){if(e===null){if(!a)throw Error(s(318));if(a=t.memoizedState,a=a!==null?a.dehydrated:null,!a)throw Error(s(317));a[rr]=t}else En(),(t.flags&128)===0&&(t.memoizedState=null),t.flags|=4;xt(t),a=!1}else Gt!==null&&(Oa(Gt),Gt=null),a=!0;if(!a)return t.flags&65536?t:null}return(t.flags&128)!==0?(t.lanes=r,t):(o=o!==null,o!==(e!==null&&e.memoizedState!==null)&&o&&(t.child.flags|=8192,(t.mode&1)!==0&&(e===null||(Ze.current&1)!==0?ct===0&&(ct=3):Fa())),t.updateQueue!==null&&(t.flags|=4),xt(t),null);case 4:return Pn(),_a(e,t),e===null&&co(t.stateNode.containerInfo),xt(t),null;case 10:return ta(t.type._context),xt(t),null;case 17:return Et(t.type)&&li(),xt(t),null;case 19:if(qe(Ze),a=t.memoizedState,a===null)return xt(t),null;if(o=(t.flags&128)!==0,d=a.rendering,d===null)if(o)ko(a,!1);else{if(ct!==0||e!==null&&(e.flags&128)!==0)for(e=t.child;e!==null;){if(d=vi(e),d!==null){for(t.flags|=128,ko(a,!1),o=d.updateQueue,o!==null&&(t.updateQueue=o,t.flags|=4),t.subtreeFlags=0,o=r,r=t.child;r!==null;)a=r,e=o,a.flags&=14680066,d=a.alternate,d===null?(a.childLanes=0,a.lanes=e,a.child=null,a.subtreeFlags=0,a.memoizedProps=null,a.memoizedState=null,a.updateQueue=null,a.dependencies=null,a.stateNode=null):(a.childLanes=d.childLanes,a.lanes=d.lanes,a.child=d.child,a.subtreeFlags=0,a.deletions=null,a.memoizedProps=d.memoizedProps,a.memoizedState=d.memoizedState,a.updateQueue=d.updateQueue,a.type=d.type,e=d.dependencies,a.dependencies=e===null?null:{lanes:e.lanes,firstContext:e.firstContext}),r=r.sibling;return Ge(Ze,Ze.current&1|2),t.child}e=e.sibling}a.tail!==null&&it()>zn&&(t.flags|=128,o=!0,ko(a,!1),t.lanes=4194304)}else{if(!o)if(e=vi(d),e!==null){if(t.flags|=128,o=!0,r=e.updateQueue,r!==null&&(t.updateQueue=r,t.flags|=4),ko(a,!0),a.tail===null&&a.tailMode==="hidden"&&!d.alternate&&!Xe)return xt(t),null}else 2*it()-a.renderingStartTime>zn&&r!==1073741824&&(t.flags|=128,o=!0,ko(a,!1),t.lanes=4194304);a.isBackwards?(d.sibling=t.child,t.child=d):(r=a.last,r!==null?r.sibling=d:t.child=d,a.last=d)}return a.tail!==null?(t=a.tail,a.rendering=t,a.tail=t.sibling,a.renderingStartTime=it(),t.sibling=null,r=Ze.current,Ge(Ze,o?r&1|2:r&1),t):(xt(t),null);case 22:case 23:return Ba(),o=t.memoizedState!==null,e!==null&&e.memoizedState!==null!==o&&(t.flags|=8192),o&&(t.mode&1)!==0?(It&1073741824)!==0&&(xt(t),t.subtreeFlags&6&&(t.flags|=8192)):xt(t),null;case 24:return null;case 25:return null}throw Error(s(156,t.tag))}function Rh(e,t){switch(Yl(t),t.tag){case 1:return Et(t.type)&&li(),e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 3:return Pn(),qe(At),qe(vt),sa(),e=t.flags,(e&65536)!==0&&(e&128)===0?(t.flags=e&-65537|128,t):null;case 5:return la(t),null;case 13:if(qe(Ze),e=t.memoizedState,e!==null&&e.dehydrated!==null){if(t.alternate===null)throw Error(s(340));En()}return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 19:return qe(Ze),null;case 4:return Pn(),null;case 10:return ta(t.type._context),null;case 22:case 23:return Ba(),null;case 24:return null;default:return null}}var Ei=!1,wt=!1,$h=typeof WeakSet=="function"?WeakSet:Set,ye=null;function $n(e,t){var r=e.ref;if(r!==null)if(typeof r=="function")try{r(null)}catch(o){rt(e,t,o)}else r.current=null}function Ta(e,t,r){try{r()}catch(o){rt(e,t,o)}}var nd=!1;function Dh(e,t){if(Bl=Qo,e=Ic(),$l(e)){if("selectionStart"in e)var r={start:e.selectionStart,end:e.selectionEnd};else e:{r=(r=e.ownerDocument)&&r.defaultView||window;var o=r.getSelection&&r.getSelection();if(o&&o.rangeCount!==0){r=o.anchorNode;var i=o.anchorOffset,a=o.focusNode;o=o.focusOffset;try{r.nodeType,a.nodeType}catch{r=null;break e}var d=0,p=-1,w=-1,O=0,Z=0,re=e,X=null;t:for(;;){for(var pe;re!==r||i!==0&&re.nodeType!==3||(p=d+i),re!==a||o!==0&&re.nodeType!==3||(w=d+o),re.nodeType===3&&(d+=re.nodeValue.length),(pe=re.firstChild)!==null;)X=re,re=pe;for(;;){if(re===e)break t;if(X===r&&++O===i&&(p=d),X===a&&++Z===o&&(w=d),(pe=re.nextSibling)!==null)break;re=X,X=re.parentNode}re=pe}r=p===-1||w===-1?null:{start:p,end:w}}else r=null}r=r||{start:0,end:0}}else r=null;for(Fl={focusedElem:e,selectionRange:r},Qo=!1,ye=t;ye!==null;)if(t=ye,e=t.child,(t.subtreeFlags&1028)!==0&&e!==null)e.return=t,ye=e;else for(;ye!==null;){t=ye;try{var we=t.alternate;if((t.flags&1024)!==0)switch(t.tag){case 0:case 11:case 15:break;case 1:if(we!==null){var ke=we.memoizedProps,lt=we.memoizedState,R=t.stateNode,k=R.getSnapshotBeforeUpdate(t.elementType===t.type?ke:Kt(t.type,ke),lt);R.__reactInternalSnapshotBeforeUpdate=k}break;case 3:var z=t.stateNode.containerInfo;z.nodeType===1?z.textContent="":z.nodeType===9&&z.documentElement&&z.removeChild(z.documentElement);break;case 5:case 6:case 4:case 17:break;default:throw Error(s(163))}}catch(ie){rt(t,t.return,ie)}if(e=t.sibling,e!==null){e.return=t.return,ye=e;break}ye=t.return}return we=nd,nd=!1,we}function So(e,t,r){var o=t.updateQueue;if(o=o!==null?o.lastEffect:null,o!==null){var i=o=o.next;do{if((i.tag&e)===e){var a=i.destroy;i.destroy=void 0,a!==void 0&&Ta(t,r,a)}i=i.next}while(i!==o)}}function _i(e,t){if(t=t.updateQueue,t=t!==null?t.lastEffect:null,t!==null){var r=t=t.next;do{if((r.tag&e)===e){var o=r.create;r.destroy=o()}r=r.next}while(r!==t)}}function ja(e){var t=e.ref;if(t!==null){var r=e.stateNode;switch(e.tag){case 5:e=r;break;default:e=r}typeof t=="function"?t(e):t.current=e}}function od(e){var t=e.alternate;t!==null&&(e.alternate=null,od(t)),e.child=null,e.deletions=null,e.sibling=null,e.tag===5&&(t=e.stateNode,t!==null&&(delete t[rr],delete t[fo],delete t[Vl],delete t[mh],delete t[gh])),e.stateNode=null,e.return=null,e.dependencies=null,e.memoizedProps=null,e.memoizedState=null,e.pendingProps=null,e.stateNode=null,e.updateQueue=null}function id(e){return e.tag===5||e.tag===3||e.tag===4}function ld(e){e:for(;;){for(;e.sibling===null;){if(e.return===null||id(e.return))return null;e=e.return}for(e.sibling.return=e.return,e=e.sibling;e.tag!==5&&e.tag!==6&&e.tag!==18;){if(e.flags&2||e.child===null||e.tag===4)continue e;e.child.return=e,e=e.child}if(!(e.flags&2))return e.stateNode}}function Pa(e,t,r){var o=e.tag;if(o===5||o===6)e=e.stateNode,t?r.nodeType===8?r.parentNode.insertBefore(e,t):r.insertBefore(e,t):(r.nodeType===8?(t=r.parentNode,t.insertBefore(e,r)):(t=r,t.appendChild(e)),r=r._reactRootContainer,r!=null||t.onclick!==null||(t.onclick=oi));else if(o!==4&&(e=e.child,e!==null))for(Pa(e,t,r),e=e.sibling;e!==null;)Pa(e,t,r),e=e.sibling}function Ra(e,t,r){var o=e.tag;if(o===5||o===6)e=e.stateNode,t?r.insertBefore(e,t):r.appendChild(e);else if(o!==4&&(e=e.child,e!==null))for(Ra(e,t,r),e=e.sibling;e!==null;)Ra(e,t,r),e=e.sibling}var ht=null,Yt=!1;function Dr(e,t,r){for(r=r.child;r!==null;)ad(e,t,r),r=r.sibling}function ad(e,t,r){if(tr&&typeof tr.onCommitFiberUnmount=="function")try{tr.onCommitFiberUnmount(Bo,r)}catch{}switch(r.tag){case 5:wt||$n(r,t);case 6:var o=ht,i=Yt;ht=null,Dr(e,t,r),ht=o,Yt=i,ht!==null&&(Yt?(e=ht,r=r.stateNode,e.nodeType===8?e.parentNode.removeChild(r):e.removeChild(r)):ht.removeChild(r.stateNode));break;case 18:ht!==null&&(Yt?(e=ht,r=r.stateNode,e.nodeType===8?Ul(e.parentNode,r):e.nodeType===1&&Ul(e,r),eo(e)):Ul(ht,r.stateNode));break;case 4:o=ht,i=Yt,ht=r.stateNode.containerInfo,Yt=!0,Dr(e,t,r),ht=o,Yt=i;break;case 0:case 11:case 14:case 15:if(!wt&&(o=r.updateQueue,o!==null&&(o=o.lastEffect,o!==null))){i=o=o.next;do{var a=i,d=a.destroy;a=a.tag,d!==void 0&&((a&2)!==0||(a&4)!==0)&&Ta(r,t,d),i=i.next}while(i!==o)}Dr(e,t,r);break;case 1:if(!wt&&($n(r,t),o=r.stateNode,typeof o.componentWillUnmount=="function"))try{o.props=r.memoizedProps,o.state=r.memoizedState,o.componentWillUnmount()}catch(p){rt(r,t,p)}Dr(e,t,r);break;case 21:Dr(e,t,r);break;case 22:r.mode&1?(wt=(o=wt)||r.memoizedState!==null,Dr(e,t,r),wt=o):Dr(e,t,r);break;default:Dr(e,t,r)}}function sd(e){var t=e.updateQueue;if(t!==null){e.updateQueue=null;var r=e.stateNode;r===null&&(r=e.stateNode=new $h),t.forEach(function(o){var i=Wh.bind(null,e,o);r.has(o)||(r.add(o),o.then(i,i))})}}function qt(e,t){var r=t.deletions;if(r!==null)for(var o=0;o<r.length;o++){var i=r[o];try{var a=e,d=t,p=d;e:for(;p!==null;){switch(p.tag){case 5:ht=p.stateNode,Yt=!1;break e;case 3:ht=p.stateNode.containerInfo,Yt=!0;break e;case 4:ht=p.stateNode.containerInfo,Yt=!0;break e}p=p.return}if(ht===null)throw Error(s(160));ad(a,d,i),ht=null,Yt=!1;var w=i.alternate;w!==null&&(w.return=null),i.return=null}catch(O){rt(i,t,O)}}if(t.subtreeFlags&12854)for(t=t.child;t!==null;)cd(t,e),t=t.sibling}function cd(e,t){var r=e.alternate,o=e.flags;switch(e.tag){case 0:case 11:case 14:case 15:if(qt(t,e),ir(e),o&4){try{So(3,e,e.return),_i(3,e)}catch(ke){rt(e,e.return,ke)}try{So(5,e,e.return)}catch(ke){rt(e,e.return,ke)}}break;case 1:qt(t,e),ir(e),o&512&&r!==null&&$n(r,r.return);break;case 5:if(qt(t,e),ir(e),o&512&&r!==null&&$n(r,r.return),e.flags&32){var i=e.stateNode;try{Hn(i,"")}catch(ke){rt(e,e.return,ke)}}if(o&4&&(i=e.stateNode,i!=null)){var a=e.memoizedProps,d=r!==null?r.memoizedProps:a,p=e.type,w=e.updateQueue;if(e.updateQueue=null,w!==null)try{p==="input"&&a.type==="radio"&&a.name!=null&&Me(i,a),sl(p,d);var O=sl(p,a);for(d=0;d<w.length;d+=2){var Z=w[d],re=w[d+1];Z==="style"?Vs(i,re):Z==="dangerouslySetInnerHTML"?Hs(i,re):Z==="children"?Hn(i,re):V(i,Z,re,O)}switch(p){case"input":Re(i,a);break;case"textarea":ve(i,a);break;case"select":var X=i._wrapperState.wasMultiple;i._wrapperState.wasMultiple=!!a.multiple;var pe=a.value;pe!=null?N(i,!!a.multiple,pe,!1):X!==!!a.multiple&&(a.defaultValue!=null?N(i,!!a.multiple,a.defaultValue,!0):N(i,!!a.multiple,a.multiple?[]:"",!1))}i[fo]=a}catch(ke){rt(e,e.return,ke)}}break;case 6:if(qt(t,e),ir(e),o&4){if(e.stateNode===null)throw Error(s(162));i=e.stateNode,a=e.memoizedProps;try{i.nodeValue=a}catch(ke){rt(e,e.return,ke)}}break;case 3:if(qt(t,e),ir(e),o&4&&r!==null&&r.memoizedState.isDehydrated)try{eo(t.containerInfo)}catch(ke){rt(e,e.return,ke)}break;case 4:qt(t,e),ir(e);break;case 13:qt(t,e),ir(e),i=e.child,i.flags&8192&&(a=i.memoizedState!==null,i.stateNode.isHidden=a,!a||i.alternate!==null&&i.alternate.memoizedState!==null||(za=it())),o&4&&sd(e);break;case 22:if(Z=r!==null&&r.memoizedState!==null,e.mode&1?(wt=(O=wt)||Z,qt(t,e),wt=O):qt(t,e),ir(e),o&8192){if(O=e.memoizedState!==null,(e.stateNode.isHidden=O)&&!Z&&(e.mode&1)!==0)for(ye=e,Z=e.child;Z!==null;){for(re=ye=Z;ye!==null;){switch(X=ye,pe=X.child,X.tag){case 0:case 11:case 14:case 15:So(4,X,X.return);break;case 1:$n(X,X.return);var we=X.stateNode;if(typeof we.componentWillUnmount=="function"){o=X,r=X.return;try{t=o,we.props=t.memoizedProps,we.state=t.memoizedState,we.componentWillUnmount()}catch(ke){rt(o,r,ke)}}break;case 5:$n(X,X.return);break;case 22:if(X.memoizedState!==null){fd(re);continue}}pe!==null?(pe.return=X,ye=pe):fd(re)}Z=Z.sibling}e:for(Z=null,re=e;;){if(re.tag===5){if(Z===null){Z=re;try{i=re.stateNode,O?(a=i.style,typeof a.setProperty=="function"?a.setProperty("display","none","important"):a.display="none"):(p=re.stateNode,w=re.memoizedProps.style,d=w!=null&&w.hasOwnProperty("display")?w.display:null,p.style.display=Us("display",d))}catch(ke){rt(e,e.return,ke)}}}else if(re.tag===6){if(Z===null)try{re.stateNode.nodeValue=O?"":re.memoizedProps}catch(ke){rt(e,e.return,ke)}}else if((re.tag!==22&&re.tag!==23||re.memoizedState===null||re===e)&&re.child!==null){re.child.return=re,re=re.child;continue}if(re===e)break e;for(;re.sibling===null;){if(re.return===null||re.return===e)break e;Z===re&&(Z=null),re=re.return}Z===re&&(Z=null),re.sibling.return=re.return,re=re.sibling}}break;case 19:qt(t,e),ir(e),o&4&&sd(e);break;case 21:break;default:qt(t,e),ir(e)}}function ir(e){var t=e.flags;if(t&2){try{e:{for(var r=e.return;r!==null;){if(id(r)){var o=r;break e}r=r.return}throw Error(s(160))}switch(o.tag){case 5:var i=o.stateNode;o.flags&32&&(Hn(i,""),o.flags&=-33);var a=ld(e);Ra(e,a,i);break;case 3:case 4:var d=o.stateNode.containerInfo,p=ld(e);Pa(e,p,d);break;default:throw Error(s(161))}}catch(w){rt(e,e.return,w)}e.flags&=-3}t&4096&&(e.flags&=-4097)}function zh(e,t,r){ye=e,ud(e)}function ud(e,t,r){for(var o=(e.mode&1)!==0;ye!==null;){var i=ye,a=i.child;if(i.tag===22&&o){var d=i.memoizedState!==null||Ei;if(!d){var p=i.alternate,w=p!==null&&p.memoizedState!==null||wt;p=Ei;var O=wt;if(Ei=d,(wt=w)&&!O)for(ye=i;ye!==null;)d=ye,w=d.child,d.tag===22&&d.memoizedState!==null?pd(i):w!==null?(w.return=d,ye=w):pd(i);for(;a!==null;)ye=a,ud(a),a=a.sibling;ye=i,Ei=p,wt=O}dd(e)}else(i.subtreeFlags&8772)!==0&&a!==null?(a.return=i,ye=a):dd(e)}}function dd(e){for(;ye!==null;){var t=ye;if((t.flags&8772)!==0){var r=t.alternate;try{if((t.flags&8772)!==0)switch(t.tag){case 0:case 11:case 15:wt||_i(5,t);break;case 1:var o=t.stateNode;if(t.flags&4&&!wt)if(r===null)o.componentDidMount();else{var i=t.elementType===t.type?r.memoizedProps:Kt(t.type,r.memoizedProps);o.componentDidUpdate(i,r.memoizedState,o.__reactInternalSnapshotBeforeUpdate)}var a=t.updateQueue;a!==null&&fu(t,a,o);break;case 3:var d=t.updateQueue;if(d!==null){if(r=null,t.child!==null)switch(t.child.tag){case 5:r=t.child.stateNode;break;case 1:r=t.child.stateNode}fu(t,d,r)}break;case 5:var p=t.stateNode;if(r===null&&t.flags&4){r=p;var w=t.memoizedProps;switch(t.type){case"button":case"input":case"select":case"textarea":w.autoFocus&&r.focus();break;case"img":w.src&&(r.src=w.src)}}break;case 6:break;case 4:break;case 12:break;case 13:if(t.memoizedState===null){var O=t.alternate;if(O!==null){var Z=O.memoizedState;if(Z!==null){var re=Z.dehydrated;re!==null&&eo(re)}}}break;case 19:case 17:case 21:case 22:case 23:case 25:break;default:throw Error(s(163))}wt||t.flags&512&&ja(t)}catch(X){rt(t,t.return,X)}}if(t===e){ye=null;break}if(r=t.sibling,r!==null){r.return=t.return,ye=r;break}ye=t.return}}function fd(e){for(;ye!==null;){var t=ye;if(t===e){ye=null;break}var r=t.sibling;if(r!==null){r.return=t.return,ye=r;break}ye=t.return}}function pd(e){for(;ye!==null;){var t=ye;try{switch(t.tag){case 0:case 11:case 15:var r=t.return;try{_i(4,t)}catch(w){rt(t,r,w)}break;case 1:var o=t.stateNode;if(typeof o.componentDidMount=="function"){var i=t.return;try{o.componentDidMount()}catch(w){rt(t,i,w)}}var a=t.return;try{ja(t)}catch(w){rt(t,a,w)}break;case 5:var d=t.return;try{ja(t)}catch(w){rt(t,d,w)}}}catch(w){rt(t,t.return,w)}if(t===e){ye=null;break}var p=t.sibling;if(p!==null){p.return=t.return,ye=p;break}ye=t.return}}var Ih=Math.ceil,Ti=J.ReactCurrentDispatcher,$a=J.ReactCurrentOwner,Bt=J.ReactCurrentBatchConfig,Oe=0,ft=null,at=null,mt=0,It=0,Dn=Tr(0),ct=0,Co=null,Zr=0,ji=0,Da=0,Ao=null,Tt=null,za=0,zn=1/0,yr=null,Pi=!1,Ia=null,zr=null,Ri=!1,Ir=null,$i=0,Eo=0,La=null,Di=-1,zi=0;function kt(){return(Oe&6)!==0?it():Di!==-1?Di:Di=it()}function Lr(e){return(e.mode&1)===0?1:(Oe&2)!==0&&mt!==0?mt&-mt:yh.transition!==null?(zi===0&&(zi=lc()),zi):(e=Ve,e!==0||(e=window.event,e=e===void 0?16:mc(e.type)),e)}function Xt(e,t,r,o){if(50<Eo)throw Eo=0,La=null,Error(s(185));Yn(e,r,o),((Oe&2)===0||e!==ft)&&(e===ft&&((Oe&2)===0&&(ji|=r),ct===4&&Mr(e,mt)),jt(e,o),r===1&&Oe===0&&(t.mode&1)===0&&(zn=it()+500,si&&Pr()))}function jt(e,t){var r=e.callbackNode;yp(e,t);var o=Ho(e,e===ft?mt:0);if(o===0)r!==null&&nc(r),e.callbackNode=null,e.callbackPriority=0;else if(t=o&-o,e.callbackPriority!==t){if(r!=null&&nc(r),t===1)e.tag===0?vh(md.bind(null,e)):eu(md.bind(null,e)),ph(function(){(Oe&6)===0&&Pr()}),r=null;else{switch(ac(o)){case 1:r=ml;break;case 4:r=oc;break;case 16:r=No;break;case 536870912:r=ic;break;default:r=No}r=Sd(r,hd.bind(null,e))}e.callbackPriority=t,e.callbackNode=r}}function hd(e,t){if(Di=-1,zi=0,(Oe&6)!==0)throw Error(s(327));var r=e.callbackNode;if(In()&&e.callbackNode!==r)return null;var o=Ho(e,e===ft?mt:0);if(o===0)return null;if((o&30)!==0||(o&e.expiredLanes)!==0||t)t=Ii(e,o);else{t=o;var i=Oe;Oe|=2;var a=vd();(ft!==e||mt!==t)&&(yr=null,zn=it()+500,en(e,t));do try{Oh();break}catch(p){gd(e,p)}while(!0);ea(),Ti.current=a,Oe=i,at!==null?t=0:(ft=null,mt=0,t=ct)}if(t!==0){if(t===2&&(i=gl(e),i!==0&&(o=i,t=Ma(e,i))),t===1)throw r=Co,en(e,0),Mr(e,o),jt(e,it()),r;if(t===6)Mr(e,o);else{if(i=e.current.alternate,(o&30)===0&&!Lh(i)&&(t=Ii(e,o),t===2&&(a=gl(e),a!==0&&(o=a,t=Ma(e,a))),t===1))throw r=Co,en(e,0),Mr(e,o),jt(e,it()),r;switch(e.finishedWork=i,e.finishedLanes=o,t){case 0:case 1:throw Error(s(345));case 2:tn(e,Tt,yr);break;case 3:if(Mr(e,o),(o&130023424)===o&&(t=za+500-it(),10<t)){if(Ho(e,0)!==0)break;if(i=e.suspendedLanes,(i&o)!==o){kt(),e.pingedLanes|=e.suspendedLanes&i;break}e.timeoutHandle=Hl(tn.bind(null,e,Tt,yr),t);break}tn(e,Tt,yr);break;case 4:if(Mr(e,o),(o&4194240)===o)break;for(t=e.eventTimes,i=-1;0<o;){var d=31-Vt(o);a=1<<d,d=t[d],d>i&&(i=d),o&=~a}if(o=i,o=it()-o,o=(120>o?120:480>o?480:1080>o?1080:1920>o?1920:3e3>o?3e3:4320>o?4320:1960*Ih(o/1960))-o,10<o){e.timeoutHandle=Hl(tn.bind(null,e,Tt,yr),o);break}tn(e,Tt,yr);break;case 5:tn(e,Tt,yr);break;default:throw Error(s(329))}}}return jt(e,it()),e.callbackNode===r?hd.bind(null,e):null}function Ma(e,t){var r=Ao;return e.current.memoizedState.isDehydrated&&(en(e,t).flags|=256),e=Ii(e,t),e!==2&&(t=Tt,Tt=r,t!==null&&Oa(t)),e}function Oa(e){Tt===null?Tt=e:Tt.push.apply(Tt,e)}function Lh(e){for(var t=e;;){if(t.flags&16384){var r=t.updateQueue;if(r!==null&&(r=r.stores,r!==null))for(var o=0;o<r.length;o++){var i=r[o],a=i.getSnapshot;i=i.value;try{if(!Qt(a(),i))return!1}catch{return!1}}}if(r=t.child,t.subtreeFlags&16384&&r!==null)r.return=t,t=r;else{if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return!0;t=t.return}t.sibling.return=t.return,t=t.sibling}}return!0}function Mr(e,t){for(t&=~Da,t&=~ji,e.suspendedLanes|=t,e.pingedLanes&=~t,e=e.expirationTimes;0<t;){var r=31-Vt(t),o=1<<r;e[r]=-1,t&=~o}}function md(e){if((Oe&6)!==0)throw Error(s(327));In();var t=Ho(e,0);if((t&1)===0)return jt(e,it()),null;var r=Ii(e,t);if(e.tag!==0&&r===2){var o=gl(e);o!==0&&(t=o,r=Ma(e,o))}if(r===1)throw r=Co,en(e,0),Mr(e,t),jt(e,it()),r;if(r===6)throw Error(s(345));return e.finishedWork=e.current.alternate,e.finishedLanes=t,tn(e,Tt,yr),jt(e,it()),null}function Na(e,t){var r=Oe;Oe|=1;try{return e(t)}finally{Oe=r,Oe===0&&(zn=it()+500,si&&Pr())}}function Jr(e){Ir!==null&&Ir.tag===0&&(Oe&6)===0&&In();var t=Oe;Oe|=1;var r=Bt.transition,o=Ve;try{if(Bt.transition=null,Ve=1,e)return e()}finally{Ve=o,Bt.transition=r,Oe=t,(Oe&6)===0&&Pr()}}function Ba(){It=Dn.current,qe(Dn)}function en(e,t){e.finishedWork=null,e.finishedLanes=0;var r=e.timeoutHandle;if(r!==-1&&(e.timeoutHandle=-1,fh(r)),at!==null)for(r=at.return;r!==null;){var o=r;switch(Yl(o),o.tag){case 1:o=o.type.childContextTypes,o!=null&&li();break;case 3:Pn(),qe(At),qe(vt),sa();break;case 5:la(o);break;case 4:Pn();break;case 13:qe(Ze);break;case 19:qe(Ze);break;case 10:ta(o.type._context);break;case 22:case 23:Ba()}r=r.return}if(ft=e,at=e=Or(e.current,null),mt=It=t,ct=0,Co=null,Da=ji=Zr=0,Tt=Ao=null,Yr!==null){for(t=0;t<Yr.length;t++)if(r=Yr[t],o=r.interleaved,o!==null){r.interleaved=null;var i=o.next,a=r.pending;if(a!==null){var d=a.next;a.next=i,o.next=d}r.pending=o}Yr=null}return e}function gd(e,t){do{var r=at;try{if(ea(),yi.current=ki,xi){for(var o=Je.memoizedState;o!==null;){var i=o.queue;i!==null&&(i.pending=null),o=o.next}xi=!1}if(Xr=0,dt=st=Je=null,yo=!1,xo=0,$a.current=null,r===null||r.return===null){ct=1,Co=t,at=null;break}e:{var a=e,d=r.return,p=r,w=t;if(t=mt,p.flags|=32768,w!==null&&typeof w=="object"&&typeof w.then=="function"){var O=w,Z=p,re=Z.tag;if((Z.mode&1)===0&&(re===0||re===11||re===15)){var X=Z.alternate;X?(Z.updateQueue=X.updateQueue,Z.memoizedState=X.memoizedState,Z.lanes=X.lanes):(Z.updateQueue=null,Z.memoizedState=null)}var pe=Fu(d);if(pe!==null){pe.flags&=-257,Wu(pe,d,p,a,t),pe.mode&1&&Bu(a,O,t),t=pe,w=O;var we=t.updateQueue;if(we===null){var ke=new Set;ke.add(w),t.updateQueue=ke}else we.add(w);break e}else{if((t&1)===0){Bu(a,O,t),Fa();break e}w=Error(s(426))}}else if(Xe&&p.mode&1){var lt=Fu(d);if(lt!==null){(lt.flags&65536)===0&&(lt.flags|=256),Wu(lt,d,p,a,t),Zl(Rn(w,p));break e}}a=w=Rn(w,p),ct!==4&&(ct=2),Ao===null?Ao=[a]:Ao.push(a),a=d;do{switch(a.tag){case 3:a.flags|=65536,t&=-t,a.lanes|=t;var R=Ou(a,w,t);du(a,R);break e;case 1:p=w;var k=a.type,z=a.stateNode;if((a.flags&128)===0&&(typeof k.getDerivedStateFromError=="function"||z!==null&&typeof z.componentDidCatch=="function"&&(zr===null||!zr.has(z)))){a.flags|=65536,t&=-t,a.lanes|=t;var ie=Nu(a,p,t);du(a,ie);break e}}a=a.return}while(a!==null)}xd(r)}catch(Se){t=Se,at===r&&r!==null&&(at=r=r.return);continue}break}while(!0)}function vd(){var e=Ti.current;return Ti.current=ki,e===null?ki:e}function Fa(){(ct===0||ct===3||ct===2)&&(ct=4),ft===null||(Zr&268435455)===0&&(ji&268435455)===0||Mr(ft,mt)}function Ii(e,t){var r=Oe;Oe|=2;var o=vd();(ft!==e||mt!==t)&&(yr=null,en(e,t));do try{Mh();break}catch(i){gd(e,i)}while(!0);if(ea(),Oe=r,Ti.current=o,at!==null)throw Error(s(261));return ft=null,mt=0,ct}function Mh(){for(;at!==null;)yd(at)}function Oh(){for(;at!==null&&!cp();)yd(at)}function yd(e){var t=kd(e.alternate,e,It);e.memoizedProps=e.pendingProps,t===null?xd(e):at=t,$a.current=null}function xd(e){var t=e;do{var r=t.alternate;if(e=t.return,(t.flags&32768)===0){if(r=Ph(r,t,It),r!==null){at=r;return}}else{if(r=Rh(r,t),r!==null){r.flags&=32767,at=r;return}if(e!==null)e.flags|=32768,e.subtreeFlags=0,e.deletions=null;else{ct=6,at=null;return}}if(t=t.sibling,t!==null){at=t;return}at=t=e}while(t!==null);ct===0&&(ct=5)}function tn(e,t,r){var o=Ve,i=Bt.transition;try{Bt.transition=null,Ve=1,Nh(e,t,r,o)}finally{Bt.transition=i,Ve=o}return null}function Nh(e,t,r,o){do In();while(Ir!==null);if((Oe&6)!==0)throw Error(s(327));r=e.finishedWork;var i=e.finishedLanes;if(r===null)return null;if(e.finishedWork=null,e.finishedLanes=0,r===e.current)throw Error(s(177));e.callbackNode=null,e.callbackPriority=0;var a=r.lanes|r.childLanes;if(xp(e,a),e===ft&&(at=ft=null,mt=0),(r.subtreeFlags&2064)===0&&(r.flags&2064)===0||Ri||(Ri=!0,Sd(No,function(){return In(),null})),a=(r.flags&15990)!==0,(r.subtreeFlags&15990)!==0||a){a=Bt.transition,Bt.transition=null;var d=Ve;Ve=1;var p=Oe;Oe|=4,$a.current=null,Dh(e,r),cd(r,e),ih(Fl),Qo=!!Bl,Fl=Bl=null,e.current=r,zh(r),up(),Oe=p,Ve=d,Bt.transition=a}else e.current=r;if(Ri&&(Ri=!1,Ir=e,$i=i),a=e.pendingLanes,a===0&&(zr=null),pp(r.stateNode),jt(e,it()),t!==null)for(o=e.onRecoverableError,r=0;r<t.length;r++)i=t[r],o(i.value,{componentStack:i.stack,digest:i.digest});if(Pi)throw Pi=!1,e=Ia,Ia=null,e;return($i&1)!==0&&e.tag!==0&&In(),a=e.pendingLanes,(a&1)!==0?e===La?Eo++:(Eo=0,La=e):Eo=0,Pr(),null}function In(){if(Ir!==null){var e=ac($i),t=Bt.transition,r=Ve;try{if(Bt.transition=null,Ve=16>e?16:e,Ir===null)var o=!1;else{if(e=Ir,Ir=null,$i=0,(Oe&6)!==0)throw Error(s(331));var i=Oe;for(Oe|=4,ye=e.current;ye!==null;){var a=ye,d=a.child;if((ye.flags&16)!==0){var p=a.deletions;if(p!==null){for(var w=0;w<p.length;w++){var O=p[w];for(ye=O;ye!==null;){var Z=ye;switch(Z.tag){case 0:case 11:case 15:So(8,Z,a)}var re=Z.child;if(re!==null)re.return=Z,ye=re;else for(;ye!==null;){Z=ye;var X=Z.sibling,pe=Z.return;if(od(Z),Z===O){ye=null;break}if(X!==null){X.return=pe,ye=X;break}ye=pe}}}var we=a.alternate;if(we!==null){var ke=we.child;if(ke!==null){we.child=null;do{var lt=ke.sibling;ke.sibling=null,ke=lt}while(ke!==null)}}ye=a}}if((a.subtreeFlags&2064)!==0&&d!==null)d.return=a,ye=d;else e:for(;ye!==null;){if(a=ye,(a.flags&2048)!==0)switch(a.tag){case 0:case 11:case 15:So(9,a,a.return)}var R=a.sibling;if(R!==null){R.return=a.return,ye=R;break e}ye=a.return}}var k=e.current;for(ye=k;ye!==null;){d=ye;var z=d.child;if((d.subtreeFlags&2064)!==0&&z!==null)z.return=d,ye=z;else e:for(d=k;ye!==null;){if(p=ye,(p.flags&2048)!==0)try{switch(p.tag){case 0:case 11:case 15:_i(9,p)}}catch(Se){rt(p,p.return,Se)}if(p===d){ye=null;break e}var ie=p.sibling;if(ie!==null){ie.return=p.return,ye=ie;break e}ye=p.return}}if(Oe=i,Pr(),tr&&typeof tr.onPostCommitFiberRoot=="function")try{tr.onPostCommitFiberRoot(Bo,e)}catch{}o=!0}return o}finally{Ve=r,Bt.transition=t}}return!1}function wd(e,t,r){t=Rn(r,t),t=Ou(e,t,1),e=$r(e,t,1),t=kt(),e!==null&&(Yn(e,1,t),jt(e,t))}function rt(e,t,r){if(e.tag===3)wd(e,e,r);else for(;t!==null;){if(t.tag===3){wd(t,e,r);break}else if(t.tag===1){var o=t.stateNode;if(typeof t.type.getDerivedStateFromError=="function"||typeof o.componentDidCatch=="function"&&(zr===null||!zr.has(o))){e=Rn(r,e),e=Nu(t,e,1),t=$r(t,e,1),e=kt(),t!==null&&(Yn(t,1,e),jt(t,e));break}}t=t.return}}function Bh(e,t,r){var o=e.pingCache;o!==null&&o.delete(t),t=kt(),e.pingedLanes|=e.suspendedLanes&r,ft===e&&(mt&r)===r&&(ct===4||ct===3&&(mt&130023424)===mt&&500>it()-za?en(e,0):Da|=r),jt(e,t)}function bd(e,t){t===0&&((e.mode&1)===0?t=1:(t=Wo,Wo<<=1,(Wo&130023424)===0&&(Wo=4194304)));var r=kt();e=mr(e,t),e!==null&&(Yn(e,t,r),jt(e,r))}function Fh(e){var t=e.memoizedState,r=0;t!==null&&(r=t.retryLane),bd(e,r)}function Wh(e,t){var r=0;switch(e.tag){case 13:var o=e.stateNode,i=e.memoizedState;i!==null&&(r=i.retryLane);break;case 19:o=e.stateNode;break;default:throw Error(s(314))}o!==null&&o.delete(t),bd(e,r)}var kd;kd=function(e,t,r){if(e!==null)if(e.memoizedProps!==t.pendingProps||At.current)_t=!0;else{if((e.lanes&r)===0&&(t.flags&128)===0)return _t=!1,jh(e,t,r);_t=(e.flags&131072)!==0}else _t=!1,Xe&&(t.flags&1048576)!==0&&tu(t,ui,t.index);switch(t.lanes=0,t.tag){case 2:var o=t.type;Ai(e,t),e=t.pendingProps;var i=Sn(t,vt.current);jn(t,r),i=da(null,t,o,e,i,r);var a=fa();return t.flags|=1,typeof i=="object"&&i!==null&&typeof i.render=="function"&&i.$$typeof===void 0?(t.tag=1,t.memoizedState=null,t.updateQueue=null,Et(o)?(a=!0,ai(t)):a=!1,t.memoizedState=i.state!==null&&i.state!==void 0?i.state:null,oa(t),i.updater=Si,t.stateNode=i,i._reactInternals=t,ya(t,o,e,r),t=ka(null,t,o,!0,a,r)):(t.tag=0,Xe&&a&&Kl(t),bt(null,t,i,r),t=t.child),t;case 16:o=t.elementType;e:{switch(Ai(e,t),e=t.pendingProps,i=o._init,o=i(o._payload),t.type=o,i=t.tag=Uh(o),e=Kt(o,e),i){case 0:t=ba(null,t,o,e,r);break e;case 1:t=Ku(null,t,o,e,r);break e;case 11:t=Hu(null,t,o,e,r);break e;case 14:t=Uu(null,t,o,Kt(o.type,e),r);break e}throw Error(s(306,o,""))}return t;case 0:return o=t.type,i=t.pendingProps,i=t.elementType===o?i:Kt(o,i),ba(e,t,o,i,r);case 1:return o=t.type,i=t.pendingProps,i=t.elementType===o?i:Kt(o,i),Ku(e,t,o,i,r);case 3:e:{if(Yu(t),e===null)throw Error(s(387));o=t.pendingProps,a=t.memoizedState,i=a.element,uu(e,t),gi(t,o,null,r);var d=t.memoizedState;if(o=d.element,a.isDehydrated)if(a={element:o,isDehydrated:!1,cache:d.cache,pendingSuspenseBoundaries:d.pendingSuspenseBoundaries,transitions:d.transitions},t.updateQueue.baseState=a,t.memoizedState=a,t.flags&256){i=Rn(Error(s(423)),t),t=qu(e,t,o,r,i);break e}else if(o!==i){i=Rn(Error(s(424)),t),t=qu(e,t,o,r,i);break e}else for(zt=_r(t.stateNode.containerInfo.firstChild),Dt=t,Xe=!0,Gt=null,r=su(t,null,o,r),t.child=r;r;)r.flags=r.flags&-3|4096,r=r.sibling;else{if(En(),o===i){t=vr(e,t,r);break e}bt(e,t,o,r)}t=t.child}return t;case 5:return pu(t),e===null&&Xl(t),o=t.type,i=t.pendingProps,a=e!==null?e.memoizedProps:null,d=i.children,Wl(o,i)?d=null:a!==null&&Wl(o,a)&&(t.flags|=32),Gu(e,t),bt(e,t,d,r),t.child;case 6:return e===null&&Xl(t),null;case 13:return Xu(e,t,r);case 4:return ia(t,t.stateNode.containerInfo),o=t.pendingProps,e===null?t.child=_n(t,null,o,r):bt(e,t,o,r),t.child;case 11:return o=t.type,i=t.pendingProps,i=t.elementType===o?i:Kt(o,i),Hu(e,t,o,i,r);case 7:return bt(e,t,t.pendingProps,r),t.child;case 8:return bt(e,t,t.pendingProps.children,r),t.child;case 12:return bt(e,t,t.pendingProps.children,r),t.child;case 10:e:{if(o=t.type._context,i=t.pendingProps,a=t.memoizedProps,d=i.value,Ge(pi,o._currentValue),o._currentValue=d,a!==null)if(Qt(a.value,d)){if(a.children===i.children&&!At.current){t=vr(e,t,r);break e}}else for(a=t.child,a!==null&&(a.return=t);a!==null;){var p=a.dependencies;if(p!==null){d=a.child;for(var w=p.firstContext;w!==null;){if(w.context===o){if(a.tag===1){w=gr(-1,r&-r),w.tag=2;var O=a.updateQueue;if(O!==null){O=O.shared;var Z=O.pending;Z===null?w.next=w:(w.next=Z.next,Z.next=w),O.pending=w}}a.lanes|=r,w=a.alternate,w!==null&&(w.lanes|=r),ra(a.return,r,t),p.lanes|=r;break}w=w.next}}else if(a.tag===10)d=a.type===t.type?null:a.child;else if(a.tag===18){if(d=a.return,d===null)throw Error(s(341));d.lanes|=r,p=d.alternate,p!==null&&(p.lanes|=r),ra(d,r,t),d=a.sibling}else d=a.child;if(d!==null)d.return=a;else for(d=a;d!==null;){if(d===t){d=null;break}if(a=d.sibling,a!==null){a.return=d.return,d=a;break}d=d.return}a=d}bt(e,t,i.children,r),t=t.child}return t;case 9:return i=t.type,o=t.pendingProps.children,jn(t,r),i=Ot(i),o=o(i),t.flags|=1,bt(e,t,o,r),t.child;case 14:return o=t.type,i=Kt(o,t.pendingProps),i=Kt(o.type,i),Uu(e,t,o,i,r);case 15:return Vu(e,t,t.type,t.pendingProps,r);case 17:return o=t.type,i=t.pendingProps,i=t.elementType===o?i:Kt(o,i),Ai(e,t),t.tag=1,Et(o)?(e=!0,ai(t)):e=!1,jn(t,r),Lu(t,o,i),ya(t,o,i,r),ka(null,t,o,!0,e,r);case 19:return Ju(e,t,r);case 22:return Qu(e,t,r)}throw Error(s(156,t.tag))};function Sd(e,t){return rc(e,t)}function Hh(e,t,r,o){this.tag=e,this.key=r,this.sibling=this.child=this.return=this.stateNode=this.type=this.elementType=null,this.index=0,this.ref=null,this.pendingProps=t,this.dependencies=this.memoizedState=this.updateQueue=this.memoizedProps=null,this.mode=o,this.subtreeFlags=this.flags=0,this.deletions=null,this.childLanes=this.lanes=0,this.alternate=null}function Ft(e,t,r,o){return new Hh(e,t,r,o)}function Wa(e){return e=e.prototype,!(!e||!e.isReactComponent)}function Uh(e){if(typeof e=="function")return Wa(e)?1:0;if(e!=null){if(e=e.$$typeof,e===W)return 11;if(e===K)return 14}return 2}function Or(e,t){var r=e.alternate;return r===null?(r=Ft(e.tag,t,e.key,e.mode),r.elementType=e.elementType,r.type=e.type,r.stateNode=e.stateNode,r.alternate=e,e.alternate=r):(r.pendingProps=t,r.type=e.type,r.flags=0,r.subtreeFlags=0,r.deletions=null),r.flags=e.flags&14680064,r.childLanes=e.childLanes,r.lanes=e.lanes,r.child=e.child,r.memoizedProps=e.memoizedProps,r.memoizedState=e.memoizedState,r.updateQueue=e.updateQueue,t=e.dependencies,r.dependencies=t===null?null:{lanes:t.lanes,firstContext:t.firstContext},r.sibling=e.sibling,r.index=e.index,r.ref=e.ref,r}function Li(e,t,r,o,i,a){var d=2;if(o=e,typeof e=="function")Wa(e)&&(d=1);else if(typeof e=="string")d=5;else e:switch(e){case ne:return rn(r.children,i,a,t);case de:d=8,i|=8;break;case Ae:return e=Ft(12,r,t,i|2),e.elementType=Ae,e.lanes=a,e;case ce:return e=Ft(13,r,t,i),e.elementType=ce,e.lanes=a,e;case xe:return e=Ft(19,r,t,i),e.elementType=xe,e.lanes=a,e;case ge:return Mi(r,i,a,t);default:if(typeof e=="object"&&e!==null)switch(e.$$typeof){case $e:d=10;break e;case De:d=9;break e;case W:d=11;break e;case K:d=14;break e;case ue:d=16,o=null;break e}throw Error(s(130,e==null?e:typeof e,""))}return t=Ft(d,r,t,i),t.elementType=e,t.type=o,t.lanes=a,t}function rn(e,t,r,o){return e=Ft(7,e,o,t),e.lanes=r,e}function Mi(e,t,r,o){return e=Ft(22,e,o,t),e.elementType=ge,e.lanes=r,e.stateNode={isHidden:!1},e}function Ha(e,t,r){return e=Ft(6,e,null,t),e.lanes=r,e}function Ua(e,t,r){return t=Ft(4,e.children!==null?e.children:[],e.key,t),t.lanes=r,t.stateNode={containerInfo:e.containerInfo,pendingChildren:null,implementation:e.implementation},t}function Vh(e,t,r,o,i){this.tag=t,this.containerInfo=e,this.finishedWork=this.pingCache=this.current=this.pendingChildren=null,this.timeoutHandle=-1,this.callbackNode=this.pendingContext=this.context=null,this.callbackPriority=0,this.eventTimes=vl(0),this.expirationTimes=vl(-1),this.entangledLanes=this.finishedLanes=this.mutableReadLanes=this.expiredLanes=this.pingedLanes=this.suspendedLanes=this.pendingLanes=0,this.entanglements=vl(0),this.identifierPrefix=o,this.onRecoverableError=i,this.mutableSourceEagerHydrationData=null}function Va(e,t,r,o,i,a,d,p,w){return e=new Vh(e,t,r,p,w),t===1?(t=1,a===!0&&(t|=8)):t=0,a=Ft(3,null,null,t),e.current=a,a.stateNode=e,a.memoizedState={element:o,isDehydrated:r,cache:null,transitions:null,pendingSuspenseBoundaries:null},oa(a),e}function Qh(e,t,r){var o=3<arguments.length&&arguments[3]!==void 0?arguments[3]:null;return{$$typeof:M,key:o==null?null:""+o,children:e,containerInfo:t,implementation:r}}function Cd(e){if(!e)return jr;e=e._reactInternals;e:{if(Ur(e)!==e||e.tag!==1)throw Error(s(170));var t=e;do{switch(t.tag){case 3:t=t.stateNode.context;break e;case 1:if(Et(t.type)){t=t.stateNode.__reactInternalMemoizedMergedChildContext;break e}}t=t.return}while(t!==null);throw Error(s(171))}if(e.tag===1){var r=e.type;if(Et(r))return Zc(e,r,t)}return t}function Ad(e,t,r,o,i,a,d,p,w){return e=Va(r,o,!0,e,i,a,d,p,w),e.context=Cd(null),r=e.current,o=kt(),i=Lr(r),a=gr(o,i),a.callback=t??null,$r(r,a,i),e.current.lanes=i,Yn(e,i,o),jt(e,o),e}function Oi(e,t,r,o){var i=t.current,a=kt(),d=Lr(i);return r=Cd(r),t.context===null?t.context=r:t.pendingContext=r,t=gr(a,d),t.payload={element:e},o=o===void 0?null:o,o!==null&&(t.callback=o),e=$r(i,t,d),e!==null&&(Xt(e,i,d,a),mi(e,i,d)),d}function Ni(e){if(e=e.current,!e.child)return null;switch(e.child.tag){case 5:return e.child.stateNode;default:return e.child.stateNode}}function Ed(e,t){if(e=e.memoizedState,e!==null&&e.dehydrated!==null){var r=e.retryLane;e.retryLane=r!==0&&r<t?r:t}}function Qa(e,t){Ed(e,t),(e=e.alternate)&&Ed(e,t)}function Gh(){return null}var _d=typeof reportError=="function"?reportError:function(e){console.error(e)};function Ga(e){this._internalRoot=e}Bi.prototype.render=Ga.prototype.render=function(e){var t=this._internalRoot;if(t===null)throw Error(s(409));Oi(e,t,null,null)},Bi.prototype.unmount=Ga.prototype.unmount=function(){var e=this._internalRoot;if(e!==null){this._internalRoot=null;var t=e.containerInfo;Jr(function(){Oi(null,e,null,null)}),t[dr]=null}};function Bi(e){this._internalRoot=e}Bi.prototype.unstable_scheduleHydration=function(e){if(e){var t=uc();e={blockedOn:null,target:e,priority:t};for(var r=0;r<Cr.length&&t!==0&&t<Cr[r].priority;r++);Cr.splice(r,0,e),r===0&&pc(e)}};function Ka(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11)}function Fi(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11&&(e.nodeType!==8||e.nodeValue!==" react-mount-point-unstable "))}function Td(){}function Kh(e,t,r,o,i){if(i){if(typeof o=="function"){var a=o;o=function(){var O=Ni(d);a.call(O)}}var d=Ad(t,o,e,0,null,!1,!1,"",Td);return e._reactRootContainer=d,e[dr]=d.current,co(e.nodeType===8?e.parentNode:e),Jr(),d}for(;i=e.lastChild;)e.removeChild(i);if(typeof o=="function"){var p=o;o=function(){var O=Ni(w);p.call(O)}}var w=Va(e,0,!1,null,null,!1,!1,"",Td);return e._reactRootContainer=w,e[dr]=w.current,co(e.nodeType===8?e.parentNode:e),Jr(function(){Oi(t,w,r,o)}),w}function Wi(e,t,r,o,i){var a=r._reactRootContainer;if(a){var d=a;if(typeof i=="function"){var p=i;i=function(){var w=Ni(d);p.call(w)}}Oi(t,d,e,i)}else d=Kh(r,t,e,i,o);return Ni(d)}sc=function(e){switch(e.tag){case 3:var t=e.stateNode;if(t.current.memoizedState.isDehydrated){var r=Kn(t.pendingLanes);r!==0&&(yl(t,r|1),jt(t,it()),(Oe&6)===0&&(zn=it()+500,Pr()))}break;case 13:Jr(function(){var o=mr(e,1);if(o!==null){var i=kt();Xt(o,e,1,i)}}),Qa(e,1)}},xl=function(e){if(e.tag===13){var t=mr(e,134217728);if(t!==null){var r=kt();Xt(t,e,134217728,r)}Qa(e,134217728)}},cc=function(e){if(e.tag===13){var t=Lr(e),r=mr(e,t);if(r!==null){var o=kt();Xt(r,e,t,o)}Qa(e,t)}},uc=function(){return Ve},dc=function(e,t){var r=Ve;try{return Ve=e,t()}finally{Ve=r}},dl=function(e,t,r){switch(t){case"input":if(Re(e,r),t=r.name,r.type==="radio"&&t!=null){for(r=e;r.parentNode;)r=r.parentNode;for(r=r.querySelectorAll("input[name="+JSON.stringify(""+t)+'][type="radio"]'),t=0;t<r.length;t++){var o=r[t];if(o!==e&&o.form===e.form){var i=ii(o);if(!i)throw Error(s(90));ae(o),Re(o,i)}}}break;case"textarea":ve(e,r);break;case"select":t=r.value,t!=null&&N(e,!!r.multiple,t,!1)}},Ys=Na,qs=Jr;var Yh={usingClientEntryPoint:!1,Events:[po,bn,ii,Gs,Ks,Na]},_o={findFiberByHostInstance:Vr,bundleType:0,version:"18.3.1",rendererPackageName:"react-dom"},qh={bundleType:_o.bundleType,version:_o.version,rendererPackageName:_o.rendererPackageName,rendererConfig:_o.rendererConfig,overrideHookState:null,overrideHookStateDeletePath:null,overrideHookStateRenamePath:null,overrideProps:null,overridePropsDeletePath:null,overridePropsRenamePath:null,setErrorHandler:null,setSuspenseHandler:null,scheduleUpdate:null,currentDispatcherRef:J.ReactCurrentDispatcher,findHostInstanceByFiber:function(e){return e=ec(e),e===null?null:e.stateNode},findFiberByHostInstance:_o.findFiberByHostInstance||Gh,findHostInstancesForRefresh:null,scheduleRefresh:null,scheduleRoot:null,setRefreshHandler:null,getCurrentFiber:null,reconcilerVersion:"18.3.1-next-f1338f8080-20240426"};if(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__<"u"){var Hi=__REACT_DEVTOOLS_GLOBAL_HOOK__;if(!Hi.isDisabled&&Hi.supportsFiber)try{Bo=Hi.inject(qh),tr=Hi}catch{}}return Pt.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=Yh,Pt.createPortal=function(e,t){var r=2<arguments.length&&arguments[2]!==void 0?arguments[2]:null;if(!Ka(t))throw Error(s(200));return Qh(e,t,null,r)},Pt.createRoot=function(e,t){if(!Ka(e))throw Error(s(299));var r=!1,o="",i=_d;return t!=null&&(t.unstable_strictMode===!0&&(r=!0),t.identifierPrefix!==void 0&&(o=t.identifierPrefix),t.onRecoverableError!==void 0&&(i=t.onRecoverableError)),t=Va(e,1,!1,null,null,r,!1,o,i),e[dr]=t.current,co(e.nodeType===8?e.parentNode:e),new Ga(t)},Pt.findDOMNode=function(e){if(e==null)return null;if(e.nodeType===1)return e;var t=e._reactInternals;if(t===void 0)throw typeof e.render=="function"?Error(s(188)):(e=Object.keys(e).join(","),Error(s(268,e)));return e=ec(t),e=e===null?null:e.stateNode,e},Pt.flushSync=function(e){return Jr(e)},Pt.hydrate=function(e,t,r){if(!Fi(t))throw Error(s(200));return Wi(null,e,t,!0,r)},Pt.hydrateRoot=function(e,t,r){if(!Ka(e))throw Error(s(405));var o=r!=null&&r.hydratedSources||null,i=!1,a="",d=_d;if(r!=null&&(r.unstable_strictMode===!0&&(i=!0),r.identifierPrefix!==void 0&&(a=r.identifierPrefix),r.onRecoverableError!==void 0&&(d=r.onRecoverableError)),t=Ad(t,null,e,1,r??null,i,!1,a,d),e[dr]=t.current,co(e),o)for(e=0;e<o.length;e++)r=o[e],i=r._getVersion,i=i(r._source),t.mutableSourceEagerHydrationData==null?t.mutableSourceEagerHydrationData=[r,i]:t.mutableSourceEagerHydrationData.push(r,i);return new Bi(t)},Pt.render=function(e,t,r){if(!Fi(t))throw Error(s(200));return Wi(null,e,t,!1,r)},Pt.unmountComponentAtNode=function(e){if(!Fi(e))throw Error(s(40));return e._reactRootContainer?(Jr(function(){Wi(null,null,e,!1,function(){e._reactRootContainer=null,e[dr]=null})}),!0):!1},Pt.unstable_batchedUpdates=Na,Pt.unstable_renderSubtreeIntoContainer=function(e,t,r,o){if(!Fi(r))throw Error(s(200));if(e==null||e._reactInternals===void 0)throw Error(s(38));return Wi(e,t,r,!1,o)},Pt.version="18.3.1-next-f1338f8080-20240426",Pt}var Ld;function kf(){if(Ld)return Xa.exports;Ld=1;function n(){if(!(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__>"u"||typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE!="function"))try{__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(n)}catch(l){console.error(l)}}return n(),Xa.exports=nm(),Xa.exports}var Md;function om(){if(Md)return Ui;Md=1;var n=kf();return Ui.createRoot=n.createRoot,Ui.hydrateRoot=n.hydrateRoot,Ui}var im=om();const lm=Ut(im);var am=`
  html,
body,
div,
span,
applet,
object,
iframe,
h1,
h2,
h3,
h4,
h5,
h6,
p,
blockquote,
pre,
a,
abbr,
acronym,
address,
big,
cite,
code,
del,
dfn,
em,
img,
ins,
kbd,
q,
s,
samp,
small,
strike,
strong,
sub,
sup,
tt,
var,
b,
u,
i,
center,
dl,
dt,
dd,
ol,
ul,
li,
fieldset,
form,
label,
legend,
table,
caption,
tbody,
tfoot,
thead,
tr,
th,
td,
article,
aside,
canvas,
details,
embed,
figure,
figcaption,
footer,
header,
hgroup,
menu,
nav,
output,
ruby,
section,
summary,
time,
mark,
audio,
video {
  margin: 0;
  padding: 0;
  border: 0;
  font-size: 100%;
  font: inherit;
  vertical-align: baseline;
}
/* HTML5 display-role reset for older browsers */
article,
aside,
details,
figcaption,
figure,
footer,
header,
hgroup,
menu,
nav,
section {
  display: block;
}
body {
  line-height: 1.5;
}
ol,
ul {
  list-style: none;
}
blockquote,
q {
  quotes: none;
}
blockquote:before,
blockquote:after,
q:before,
q:after {
  content: "";
  content: none;
}
table {
  border-collapse: collapse;
  border-spacing: 0;
}
a {
  color: inherit;
  text-decoration: none;
}
ul,
li {
  list-style-type: none;
}
button {
  outline: none;
  -webkit-tap-highlight-color: rgba(0, 0, 0, 0);
}
body {
  margin: 0;
  padding: 0;
  font-family: sans-serif;
  color: black;
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, "Courier New",
    monospace;
}

input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type="number"] {
  -moz-appearance: textfield;
}

`,es={exports:{}},Qe={};/**
 * @license React
 * react-is.production.js
 *
 * Copyright (c) Meta Platforms, Inc. and affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var Od;function sm(){if(Od)return Qe;Od=1;var n=Symbol.for("react.transitional.element"),l=Symbol.for("react.portal"),s=Symbol.for("react.fragment"),u=Symbol.for("react.strict_mode"),f=Symbol.for("react.profiler"),g=Symbol.for("react.consumer"),h=Symbol.for("react.context"),v=Symbol.for("react.forward_ref"),x=Symbol.for("react.suspense"),E=Symbol.for("react.suspense_list"),T=Symbol.for("react.memo"),A=Symbol.for("react.lazy"),P=Symbol.for("react.view_transition"),U=Symbol.for("react.client.reference");function j(y){if(typeof y=="object"&&y!==null){var $=y.$$typeof;switch($){case n:switch(y=y.type,y){case s:case f:case u:case x:case E:case P:return y;default:switch(y=y&&y.$$typeof,y){case h:case v:case A:case T:return y;case g:return y;default:return $}}case l:return $}}}return Qe.ContextConsumer=g,Qe.ContextProvider=h,Qe.Element=n,Qe.ForwardRef=v,Qe.Fragment=s,Qe.Lazy=A,Qe.Memo=T,Qe.Portal=l,Qe.Profiler=f,Qe.StrictMode=u,Qe.Suspense=x,Qe.SuspenseList=E,Qe.isContextConsumer=function(y){return j(y)===g},Qe.isContextProvider=function(y){return j(y)===h},Qe.isElement=function(y){return typeof y=="object"&&y!==null&&y.$$typeof===n},Qe.isForwardRef=function(y){return j(y)===v},Qe.isFragment=function(y){return j(y)===s},Qe.isLazy=function(y){return j(y)===A},Qe.isMemo=function(y){return j(y)===T},Qe.isPortal=function(y){return j(y)===l},Qe.isProfiler=function(y){return j(y)===f},Qe.isStrictMode=function(y){return j(y)===u},Qe.isSuspense=function(y){return j(y)===x},Qe.isSuspenseList=function(y){return j(y)===E},Qe.isValidElementType=function(y){return typeof y=="string"||typeof y=="function"||y===s||y===f||y===u||y===x||y===E||y===P||typeof y=="object"&&y!==null&&(y.$$typeof===A||y.$$typeof===T||y.$$typeof===h||y.$$typeof===g||y.$$typeof===v||y.$$typeof===U||y.getModuleId!==void 0)},Qe.typeOf=j,Qe}var Nd;function cm(){return Nd||(Nd=1,es.exports=sm()),es.exports}var Sf=cm();function um(n){function l(C,I,D,q,b){for(var oe=0,Q=0,_=0,le=0,fe,ae,je=0,ze=0,Ce,Me=Ce=fe=0,Re=0,Ne=0,ot=0,Fe=0,N=D.length,te=N-1,Ee,ve="",Ue="",ur="",Hr="",Rt;Re<N;){if(ae=D.charCodeAt(Re),Re===te&&Q+le+_+oe!==0&&(Q!==0&&(ae=Q===47?10:47),le=_=oe=0,N++,te++),Q+le+_+oe===0){if(Re===te&&(0<Ne&&(ve=ve.replace(P,"")),0<ve.trim().length)){switch(ae){case 32:case 9:case 59:case 13:case 10:break;default:ve+=D.charAt(Re)}ae=59}switch(ae){case 123:for(ve=ve.trim(),fe=ve.charCodeAt(0),Ce=1,Fe=++Re;Re<N;){switch(ae=D.charCodeAt(Re)){case 123:Ce++;break;case 125:Ce--;break;case 47:switch(ae=D.charCodeAt(Re+1)){case 42:case 47:e:{for(Me=Re+1;Me<te;++Me)switch(D.charCodeAt(Me)){case 47:if(ae===42&&D.charCodeAt(Me-1)===42&&Re+2!==Me){Re=Me+1;break e}break;case 10:if(ae===47){Re=Me+1;break e}}Re=Me}}break;case 91:ae++;case 40:ae++;case 34:case 39:for(;Re++<te&&D.charCodeAt(Re)!==ae;);}if(Ce===0)break;Re++}switch(Ce=D.substring(Fe,Re),fe===0&&(fe=(ve=ve.replace(A,"").trim()).charCodeAt(0)),fe){case 64:switch(0<Ne&&(ve=ve.replace(P,"")),ae=ve.charCodeAt(1),ae){case 100:case 109:case 115:case 45:Ne=I;break;default:Ne=ue}if(Ce=l(I,Ne,Ce,ae,b+1),Fe=Ce.length,0<F&&(Ne=s(ue,ve,ot),Rt=v(3,Ce,Ne,I,ce,W,Fe,ae,b,q),ve=Ne.join(""),Rt!==void 0&&(Fe=(Ce=Rt.trim()).length)===0&&(ae=0,Ce="")),0<Fe)switch(ae){case 115:ve=ve.replace(M,h);case 100:case 109:case 45:Ce=ve+"{"+Ce+"}";break;case 107:ve=ve.replace(H,"$1 $2"),Ce=ve+"{"+Ce+"}",Ce=K===1||K===2&&g("@"+Ce,3)?"@-webkit-"+Ce+"@"+Ce:"@"+Ce;break;default:Ce=ve+Ce,q===112&&(Ce=(Ue+=Ce,""))}else Ce="";break;default:Ce=l(I,s(I,ve,ot),Ce,q,b+1)}ur+=Ce,Ce=ot=Ne=Me=fe=0,ve="",ae=D.charCodeAt(++Re);break;case 125:case 59:if(ve=(0<Ne?ve.replace(P,""):ve).trim(),1<(Fe=ve.length))switch(Me===0&&(fe=ve.charCodeAt(0),fe===45||96<fe&&123>fe)&&(Fe=(ve=ve.replace(" ",":")).length),0<F&&(Rt=v(1,ve,I,C,ce,W,Ue.length,q,b,q))!==void 0&&(Fe=(ve=Rt.trim()).length)===0&&(ve="\0\0"),fe=ve.charCodeAt(0),ae=ve.charCodeAt(1),fe){case 0:break;case 64:if(ae===105||ae===99){Hr+=ve+D.charAt(Re);break}default:ve.charCodeAt(Fe-1)!==58&&(Ue+=f(ve,fe,ae,ve.charCodeAt(2)))}ot=Ne=Me=fe=0,ve="",ae=D.charCodeAt(++Re)}}switch(ae){case 13:case 10:Q===47?Q=0:1+fe===0&&q!==107&&0<ve.length&&(Ne=1,ve+="\0"),0<F*G&&v(0,ve,I,C,ce,W,Ue.length,q,b,q),W=1,ce++;break;case 59:case 125:if(Q+le+_+oe===0){W++;break}default:switch(W++,Ee=D.charAt(Re),ae){case 9:case 32:if(le+oe+Q===0)switch(je){case 44:case 58:case 9:case 32:Ee="";break;default:ae!==32&&(Ee=" ")}break;case 0:Ee="\\0";break;case 12:Ee="\\f";break;case 11:Ee="\\v";break;case 38:le+Q+oe===0&&(Ne=ot=1,Ee="\f"+Ee);break;case 108:if(le+Q+oe+xe===0&&0<Me)switch(Re-Me){case 2:je===112&&D.charCodeAt(Re-3)===58&&(xe=je);case 8:ze===111&&(xe=ze)}break;case 58:le+Q+oe===0&&(Me=Re);break;case 44:Q+_+le+oe===0&&(Ne=1,Ee+="\r");break;case 34:case 39:Q===0&&(le=le===ae?0:le===0?ae:le);break;case 91:le+Q+_===0&&oe++;break;case 93:le+Q+_===0&&oe--;break;case 41:le+Q+oe===0&&_--;break;case 40:if(le+Q+oe===0){if(fe===0)switch(2*je+3*ze){case 533:break;default:fe=1}_++}break;case 64:Q+_+le+oe+Me+Ce===0&&(Ce=1);break;case 42:case 47:if(!(0<le+oe+_))switch(Q){case 0:switch(2*ae+3*D.charCodeAt(Re+1)){case 235:Q=47;break;case 220:Fe=Re,Q=42}break;case 42:ae===47&&je===42&&Fe+2!==Re&&(D.charCodeAt(Fe+2)===33&&(Ue+=D.substring(Fe,Re+1)),Ee="",Q=0)}}Q===0&&(ve+=Ee)}ze=je,je=ae,Re++}if(Fe=Ue.length,0<Fe){if(Ne=I,0<F&&(Rt=v(2,Ue,Ne,C,ce,W,Fe,q,b,q),Rt!==void 0&&(Ue=Rt).length===0))return Hr+Ue+ur;if(Ue=Ne.join(",")+"{"+Ue+"}",K*xe!==0){switch(K!==2||g(Ue,2)||(xe=0),xe){case 111:Ue=Ue.replace(J,":-moz-$1")+Ue;break;case 112:Ue=Ue.replace(V,"::-webkit-input-$1")+Ue.replace(V,"::-moz-$1")+Ue.replace(V,":-ms-input-$1")+Ue}xe=0}}return Hr+Ue+ur}function s(C,I,D){var q=I.trim().split($);I=q;var b=q.length,oe=C.length;switch(oe){case 0:case 1:var Q=0;for(C=oe===0?"":C[0]+" ";Q<b;++Q)I[Q]=u(C,I[Q],D).trim();break;default:var _=Q=0;for(I=[];Q<b;++Q)for(var le=0;le<oe;++le)I[_++]=u(C[le]+" ",q[Q],D).trim()}return I}function u(C,I,D){var q=I.charCodeAt(0);switch(33>q&&(q=(I=I.trim()).charCodeAt(0)),q){case 38:return I.replace(Y,"$1"+C.trim());case 58:return C.trim()+I.replace(Y,"$1"+C.trim());default:if(0<1*D&&0<I.indexOf("\f"))return I.replace(Y,(C.charCodeAt(0)===58?"":"$1")+C.trim())}return C+I}function f(C,I,D,q){var b=C+";",oe=2*I+3*D+4*q;if(oe===944){C=b.indexOf(":",9)+1;var Q=b.substring(C,b.length-1).trim();return Q=b.substring(0,C).trim()+Q+";",K===1||K===2&&g(Q,1)?"-webkit-"+Q+Q:Q}if(K===0||K===2&&!g(b,1))return b;switch(oe){case 1015:return b.charCodeAt(10)===97?"-webkit-"+b+b:b;case 951:return b.charCodeAt(3)===116?"-webkit-"+b+b:b;case 963:return b.charCodeAt(5)===110?"-webkit-"+b+b:b;case 1009:if(b.charCodeAt(4)!==100)break;case 969:case 942:return"-webkit-"+b+b;case 978:return"-webkit-"+b+"-moz-"+b+b;case 1019:case 983:return"-webkit-"+b+"-moz-"+b+"-ms-"+b+b;case 883:if(b.charCodeAt(8)===45)return"-webkit-"+b+b;if(0<b.indexOf("image-set(",11))return b.replace(De,"$1-webkit-$2")+b;break;case 932:if(b.charCodeAt(4)===45)switch(b.charCodeAt(5)){case 103:return"-webkit-box-"+b.replace("-grow","")+"-webkit-"+b+"-ms-"+b.replace("grow","positive")+b;case 115:return"-webkit-"+b+"-ms-"+b.replace("shrink","negative")+b;case 98:return"-webkit-"+b+"-ms-"+b.replace("basis","preferred-size")+b}return"-webkit-"+b+"-ms-"+b+b;case 964:return"-webkit-"+b+"-ms-flex-"+b+b;case 1023:if(b.charCodeAt(8)!==99)break;return Q=b.substring(b.indexOf(":",15)).replace("flex-","").replace("space-between","justify"),"-webkit-box-pack"+Q+"-webkit-"+b+"-ms-flex-pack"+Q+b;case 1005:return j.test(b)?b.replace(U,":-webkit-")+b.replace(U,":-moz-")+b:b;case 1e3:switch(Q=b.substring(13).trim(),I=Q.indexOf("-")+1,Q.charCodeAt(0)+Q.charCodeAt(I)){case 226:Q=b.replace(L,"tb");break;case 232:Q=b.replace(L,"tb-rl");break;case 220:Q=b.replace(L,"lr");break;default:return b}return"-webkit-"+b+"-ms-"+Q+b;case 1017:if(b.indexOf("sticky",9)===-1)break;case 975:switch(I=(b=C).length-10,Q=(b.charCodeAt(I)===33?b.substring(0,I):b).substring(C.indexOf(":",7)+1).trim(),oe=Q.charCodeAt(0)+(Q.charCodeAt(7)|0)){case 203:if(111>Q.charCodeAt(8))break;case 115:b=b.replace(Q,"-webkit-"+Q)+";"+b;break;case 207:case 102:b=b.replace(Q,"-webkit-"+(102<oe?"inline-":"")+"box")+";"+b.replace(Q,"-webkit-"+Q)+";"+b.replace(Q,"-ms-"+Q+"box")+";"+b}return b+";";case 938:if(b.charCodeAt(5)===45)switch(b.charCodeAt(6)){case 105:return Q=b.replace("-items",""),"-webkit-"+b+"-webkit-box-"+Q+"-ms-flex-"+Q+b;case 115:return"-webkit-"+b+"-ms-flex-item-"+b.replace(de,"")+b;default:return"-webkit-"+b+"-ms-flex-line-pack"+b.replace("align-content","").replace(de,"")+b}break;case 973:case 989:if(b.charCodeAt(3)!==45||b.charCodeAt(4)===122)break;case 931:case 953:if($e.test(C)===!0)return(Q=C.substring(C.indexOf(":")+1)).charCodeAt(0)===115?f(C.replace("stretch","fill-available"),I,D,q).replace(":fill-available",":stretch"):b.replace(Q,"-webkit-"+Q)+b.replace(Q,"-moz-"+Q.replace("fill-",""))+b;break;case 962:if(b="-webkit-"+b+(b.charCodeAt(5)===102?"-ms-"+b:"")+b,D+q===211&&b.charCodeAt(13)===105&&0<b.indexOf("transform",10))return b.substring(0,b.indexOf(";",27)+1).replace(y,"$1-webkit-$2")+b}return b}function g(C,I){var D=C.indexOf(I===1?":":"{"),q=C.substring(0,I!==3?D:10);return D=C.substring(D+1,C.length-1),se(I!==2?q:q.replace(Ae,"$1"),D,I)}function h(C,I){var D=f(I,I.charCodeAt(0),I.charCodeAt(1),I.charCodeAt(2));return D!==I+";"?D.replace(ne," or ($1)").substring(4):"("+I+")"}function v(C,I,D,q,b,oe,Q,_,le,fe){for(var ae=0,je=I,ze;ae<F;++ae)switch(ze=ge[ae].call(T,C,je,D,q,b,oe,Q,_,le,fe)){case void 0:case!1:case!0:case null:break;default:je=ze}if(je!==I)return je}function x(C){switch(C){case void 0:case null:F=ge.length=0;break;default:if(typeof C=="function")ge[F++]=C;else if(typeof C=="object")for(var I=0,D=C.length;I<D;++I)x(C[I]);else G=!!C|0}return x}function E(C){return C=C.prefix,C!==void 0&&(se=null,C?typeof C!="function"?K=1:(K=2,se=C):K=0),E}function T(C,I){var D=C;if(33>D.charCodeAt(0)&&(D=D.trim()),S=D,D=[S],0<F){var q=v(-1,I,D,D,ce,W,0,0,0,0);q!==void 0&&typeof q=="string"&&(I=q)}var b=l(ue,D,I,0,0);return 0<F&&(q=v(-2,b,D,D,ce,W,b.length,0,0,0),q!==void 0&&(b=q)),S="",xe=0,W=ce=1,b}var A=/^\0+/g,P=/[\0\r\f]/g,U=/: */g,j=/zoo|gra/,y=/([,: ])(transform)/g,$=/,\r+?/g,Y=/([\t\r\n ])*\f?&/g,H=/@(k\w+)\s*(\S*)\s*/,V=/::(place)/g,J=/:(read-only)/g,L=/[svh]\w+-[tblr]{2}/,M=/\(\s*(.*)\s*\)/g,ne=/([\s\S]*?);/g,de=/-self|flex-/g,Ae=/[^]*?(:[rp][el]a[\w-]+)[^]*/,$e=/stretch|:\s*\w+\-(?:conte|avail)/,De=/([^-])(image-set\()/,W=1,ce=1,xe=0,K=1,ue=[],ge=[],F=0,se=null,G=0,S="";return T.use=x,T.set=E,n!==void 0&&E(n),T}var dm={animationIterationCount:1,borderImageOutset:1,borderImageSlice:1,borderImageWidth:1,boxFlex:1,boxFlexGroup:1,boxOrdinalGroup:1,columnCount:1,columns:1,flex:1,flexGrow:1,flexPositive:1,flexShrink:1,flexNegative:1,flexOrder:1,gridRow:1,gridRowEnd:1,gridRowSpan:1,gridRowStart:1,gridColumn:1,gridColumnEnd:1,gridColumnSpan:1,gridColumnStart:1,msGridRow:1,msGridRowSpan:1,msGridColumn:1,msGridColumnSpan:1,fontWeight:1,lineHeight:1,opacity:1,order:1,orphans:1,tabSize:1,widows:1,zIndex:1,zoom:1,WebkitLineClamp:1,fillOpacity:1,floodOpacity:1,stopOpacity:1,strokeDasharray:1,strokeDashoffset:1,strokeMiterlimit:1,strokeOpacity:1,strokeWidth:1};function fm(n){var l=Object.create(null);return function(s){return l[s]===void 0&&(l[s]=n(s)),l[s]}}var pm=/^((children|dangerouslySetInnerHTML|key|ref|autoFocus|defaultValue|defaultChecked|innerHTML|suppressContentEditableWarning|suppressHydrationWarning|valueLink|abbr|accept|acceptCharset|accessKey|action|allow|allowUserMedia|allowPaymentRequest|allowFullScreen|allowTransparency|alt|async|autoComplete|autoPlay|capture|cellPadding|cellSpacing|challenge|charSet|checked|cite|classID|className|cols|colSpan|content|contentEditable|contextMenu|controls|controlsList|coords|crossOrigin|data|dateTime|decoding|default|defer|dir|disabled|disablePictureInPicture|disableRemotePlayback|download|draggable|encType|enterKeyHint|fetchpriority|fetchPriority|form|formAction|formEncType|formMethod|formNoValidate|formTarget|frameBorder|headers|height|hidden|high|href|hrefLang|htmlFor|httpEquiv|id|inputMode|integrity|is|keyParams|keyType|kind|label|lang|list|loading|loop|low|marginHeight|marginWidth|max|maxLength|media|mediaGroup|method|min|minLength|multiple|muted|name|nonce|noValidate|open|optimum|pattern|placeholder|playsInline|popover|popoverTarget|popoverTargetAction|poster|preload|profile|radioGroup|readOnly|referrerPolicy|rel|required|reversed|role|rows|rowSpan|sandbox|scope|scoped|scrolling|seamless|selected|shape|size|sizes|slot|span|spellCheck|src|srcDoc|srcLang|srcSet|start|step|style|summary|tabIndex|target|title|translate|type|useMap|value|width|wmode|wrap|about|datatype|inlist|prefix|property|resource|typeof|vocab|autoCapitalize|autoCorrect|autoSave|color|incremental|fallback|inert|itemProp|itemScope|itemType|itemID|itemRef|on|option|results|security|unselectable|accentHeight|accumulate|additive|alignmentBaseline|allowReorder|alphabetic|amplitude|arabicForm|ascent|attributeName|attributeType|autoReverse|azimuth|baseFrequency|baselineShift|baseProfile|bbox|begin|bias|by|calcMode|capHeight|clip|clipPathUnits|clipPath|clipRule|colorInterpolation|colorInterpolationFilters|colorProfile|colorRendering|contentScriptType|contentStyleType|cursor|cx|cy|d|decelerate|descent|diffuseConstant|direction|display|divisor|dominantBaseline|dur|dx|dy|edgeMode|elevation|enableBackground|end|exponent|externalResourcesRequired|fill|fillOpacity|fillRule|filter|filterRes|filterUnits|floodColor|floodOpacity|focusable|fontFamily|fontSize|fontSizeAdjust|fontStretch|fontStyle|fontVariant|fontWeight|format|from|fr|fx|fy|g1|g2|glyphName|glyphOrientationHorizontal|glyphOrientationVertical|glyphRef|gradientTransform|gradientUnits|hanging|horizAdvX|horizOriginX|ideographic|imageRendering|in|in2|intercept|k|k1|k2|k3|k4|kernelMatrix|kernelUnitLength|kerning|keyPoints|keySplines|keyTimes|lengthAdjust|letterSpacing|lightingColor|limitingConeAngle|local|markerEnd|markerMid|markerStart|markerHeight|markerUnits|markerWidth|mask|maskContentUnits|maskUnits|mathematical|mode|numOctaves|offset|opacity|operator|order|orient|orientation|origin|overflow|overlinePosition|overlineThickness|panose1|paintOrder|pathLength|patternContentUnits|patternTransform|patternUnits|pointerEvents|points|pointsAtX|pointsAtY|pointsAtZ|preserveAlpha|preserveAspectRatio|primitiveUnits|r|radius|refX|refY|renderingIntent|repeatCount|repeatDur|requiredExtensions|requiredFeatures|restart|result|rotate|rx|ry|scale|seed|shapeRendering|slope|spacing|specularConstant|specularExponent|speed|spreadMethod|startOffset|stdDeviation|stemh|stemv|stitchTiles|stopColor|stopOpacity|strikethroughPosition|strikethroughThickness|string|stroke|strokeDasharray|strokeDashoffset|strokeLinecap|strokeLinejoin|strokeMiterlimit|strokeOpacity|strokeWidth|surfaceScale|systemLanguage|tableValues|targetX|targetY|textAnchor|textDecoration|textRendering|textLength|to|transform|u1|u2|underlinePosition|underlineThickness|unicode|unicodeBidi|unicodeRange|unitsPerEm|vAlphabetic|vHanging|vIdeographic|vMathematical|values|vectorEffect|version|vertAdvY|vertOriginX|vertOriginY|viewBox|viewTarget|visibility|widths|wordSpacing|writingMode|x|xHeight|x1|x2|xChannelSelector|xlinkActuate|xlinkArcrole|xlinkHref|xlinkRole|xlinkShow|xlinkTitle|xlinkType|xmlBase|xmlns|xmlnsXlink|xmlLang|xmlSpace|y|y1|y2|yChannelSelector|z|zoomAndPan|for|class|autofocus)|(([Dd][Aa][Tt][Aa]|[Aa][Rr][Ii][Aa]|x)-.*))$/,Bd=fm(function(n){return pm.test(n)||n.charCodeAt(0)===111&&n.charCodeAt(1)===110&&n.charCodeAt(2)<91}),ts={exports:{}},We={};/** @license React v16.13.1
 * react-is.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var Fd;function hm(){if(Fd)return We;Fd=1;var n=typeof Symbol=="function"&&Symbol.for,l=n?Symbol.for("react.element"):60103,s=n?Symbol.for("react.portal"):60106,u=n?Symbol.for("react.fragment"):60107,f=n?Symbol.for("react.strict_mode"):60108,g=n?Symbol.for("react.profiler"):60114,h=n?Symbol.for("react.provider"):60109,v=n?Symbol.for("react.context"):60110,x=n?Symbol.for("react.async_mode"):60111,E=n?Symbol.for("react.concurrent_mode"):60111,T=n?Symbol.for("react.forward_ref"):60112,A=n?Symbol.for("react.suspense"):60113,P=n?Symbol.for("react.suspense_list"):60120,U=n?Symbol.for("react.memo"):60115,j=n?Symbol.for("react.lazy"):60116,y=n?Symbol.for("react.block"):60121,$=n?Symbol.for("react.fundamental"):60117,Y=n?Symbol.for("react.responder"):60118,H=n?Symbol.for("react.scope"):60119;function V(L){if(typeof L=="object"&&L!==null){var M=L.$$typeof;switch(M){case l:switch(L=L.type,L){case x:case E:case u:case g:case f:case A:return L;default:switch(L=L&&L.$$typeof,L){case v:case T:case j:case U:case h:return L;default:return M}}case s:return M}}}function J(L){return V(L)===E}return We.AsyncMode=x,We.ConcurrentMode=E,We.ContextConsumer=v,We.ContextProvider=h,We.Element=l,We.ForwardRef=T,We.Fragment=u,We.Lazy=j,We.Memo=U,We.Portal=s,We.Profiler=g,We.StrictMode=f,We.Suspense=A,We.isAsyncMode=function(L){return J(L)||V(L)===x},We.isConcurrentMode=J,We.isContextConsumer=function(L){return V(L)===v},We.isContextProvider=function(L){return V(L)===h},We.isElement=function(L){return typeof L=="object"&&L!==null&&L.$$typeof===l},We.isForwardRef=function(L){return V(L)===T},We.isFragment=function(L){return V(L)===u},We.isLazy=function(L){return V(L)===j},We.isMemo=function(L){return V(L)===U},We.isPortal=function(L){return V(L)===s},We.isProfiler=function(L){return V(L)===g},We.isStrictMode=function(L){return V(L)===f},We.isSuspense=function(L){return V(L)===A},We.isValidElementType=function(L){return typeof L=="string"||typeof L=="function"||L===u||L===E||L===g||L===f||L===A||L===P||typeof L=="object"&&L!==null&&(L.$$typeof===j||L.$$typeof===U||L.$$typeof===h||L.$$typeof===v||L.$$typeof===T||L.$$typeof===$||L.$$typeof===Y||L.$$typeof===H||L.$$typeof===y)},We.typeOf=V,We}var Wd;function mm(){return Wd||(Wd=1,ts.exports=hm()),ts.exports}var rs,Hd;function gm(){if(Hd)return rs;Hd=1;var n=mm(),l={childContextTypes:!0,contextType:!0,contextTypes:!0,defaultProps:!0,displayName:!0,getDefaultProps:!0,getDerivedStateFromError:!0,getDerivedStateFromProps:!0,mixins:!0,propTypes:!0,type:!0},s={name:!0,length:!0,prototype:!0,caller:!0,callee:!0,arguments:!0,arity:!0},u={$$typeof:!0,render:!0,defaultProps:!0,displayName:!0,propTypes:!0},f={$$typeof:!0,compare:!0,defaultProps:!0,displayName:!0,propTypes:!0,type:!0},g={};g[n.ForwardRef]=u,g[n.Memo]=f;function h(j){return n.isMemo(j)?f:g[j.$$typeof]||l}var v=Object.defineProperty,x=Object.getOwnPropertyNames,E=Object.getOwnPropertySymbols,T=Object.getOwnPropertyDescriptor,A=Object.getPrototypeOf,P=Object.prototype;function U(j,y,$){if(typeof y!="string"){if(P){var Y=A(y);Y&&Y!==P&&U(j,Y,$)}var H=x(y);E&&(H=H.concat(E(y)));for(var V=h(j),J=h(y),L=0;L<H.length;++L){var M=H[L];if(!s[M]&&!($&&$[M])&&!(J&&J[M])&&!(V&&V[M])){var ne=T(y,M);try{v(j,M,ne)}catch{}}}}return j}return rs=U,rs}var vm=gm();const ym=Ut(vm);var Wt={};function Zt(){return(Zt=Object.assign||function(n){for(var l=1;l<arguments.length;l++){var s=arguments[l];for(var u in s)Object.prototype.hasOwnProperty.call(s,u)&&(n[u]=s[u])}return n}).apply(this,arguments)}var Ud=function(n,l){for(var s=[n[0]],u=0,f=l.length;u<f;u+=1)s.push(l[u],n[u+1]);return s},bs=function(n){return n!==null&&typeof n=="object"&&(n.toString?n.toString():Object.prototype.toString.call(n))==="[object Object]"&&!Sf.typeOf(n)},Ji=Object.freeze([]),Br=Object.freeze({});function Bn(n){return typeof n=="function"}function Vd(n){return n.displayName||n.name||"Component"}function $s(n){return n&&typeof n.styledComponentId=="string"}var Fn=typeof process<"u"&&Wt!==void 0&&(Wt.REACT_APP_SC_ATTR||Wt.SC_ATTR)||"data-styled",Ds=typeof window<"u"&&"HTMLElement"in window,xm=!!(typeof SC_DISABLE_SPEEDY=="boolean"?SC_DISABLE_SPEEDY:typeof process<"u"&&Wt!==void 0&&(Wt.REACT_APP_SC_DISABLE_SPEEDY!==void 0&&Wt.REACT_APP_SC_DISABLE_SPEEDY!==""?Wt.REACT_APP_SC_DISABLE_SPEEDY!=="false"&&Wt.REACT_APP_SC_DISABLE_SPEEDY:Wt.SC_DISABLE_SPEEDY!==void 0&&Wt.SC_DISABLE_SPEEDY!==""&&Wt.SC_DISABLE_SPEEDY!=="false"&&Wt.SC_DISABLE_SPEEDY)),wm={};function sn(n){for(var l=arguments.length,s=new Array(l>1?l-1:0),u=1;u<l;u++)s[u-1]=arguments[u];throw new Error("An error occurred. See https://git.io/JUIaE#"+n+" for more information."+(s.length>0?" Args: "+s.join(", "):""))}var bm=(function(){function n(s){this.groupSizes=new Uint32Array(512),this.length=512,this.tag=s}var l=n.prototype;return l.indexOfGroup=function(s){for(var u=0,f=0;f<s;f++)u+=this.groupSizes[f];return u},l.insertRules=function(s,u){if(s>=this.groupSizes.length){for(var f=this.groupSizes,g=f.length,h=g;s>=h;)(h<<=1)<0&&sn(16,""+s);this.groupSizes=new Uint32Array(h),this.groupSizes.set(f),this.length=h;for(var v=g;v<h;v++)this.groupSizes[v]=0}for(var x=this.indexOfGroup(s+1),E=0,T=u.length;E<T;E++)this.tag.insertRule(x,u[E])&&(this.groupSizes[s]++,x++)},l.clearGroup=function(s){if(s<this.length){var u=this.groupSizes[s],f=this.indexOfGroup(s),g=f+u;this.groupSizes[s]=0;for(var h=f;h<g;h++)this.tag.deleteRule(f)}},l.getGroup=function(s){var u="";if(s>=this.length||this.groupSizes[s]===0)return u;for(var f=this.groupSizes[s],g=this.indexOfGroup(s),h=g+f,v=g;v<h;v++)u+=this.tag.getRule(v)+`/*!sc*/
`;return u},n})(),Zi=new Map,el=new Map,jo=1,Vi=function(n){if(Zi.has(n))return Zi.get(n);for(;el.has(jo);)jo++;var l=jo++;return Zi.set(n,l),el.set(l,n),l},km=function(n){return el.get(n)},Sm=function(n,l){l>=jo&&(jo=l+1),Zi.set(n,l),el.set(l,n)},Cm="style["+Fn+'][data-styled-version="5.3.11"]',Am=new RegExp("^"+Fn+'\\.g(\\d+)\\[id="([\\w\\d-]+)"\\].*?"([^"]*)'),Em=function(n,l,s){for(var u,f=s.split(","),g=0,h=f.length;g<h;g++)(u=f[g])&&n.registerName(l,u)},_m=function(n,l){for(var s=(l.textContent||"").split(`/*!sc*/
`),u=[],f=0,g=s.length;f<g;f++){var h=s[f].trim();if(h){var v=h.match(Am);if(v){var x=0|parseInt(v[1],10),E=v[2];x!==0&&(Sm(E,x),Em(n,E,v[3]),n.getTag().insertRules(x,u)),u.length=0}else u.push(h)}}},Tm=function(){return typeof __webpack_nonce__<"u"?__webpack_nonce__:null},Cf=function(n){var l=document.head,s=n||l,u=document.createElement("style"),f=(function(v){for(var x=v.childNodes,E=x.length;E>=0;E--){var T=x[E];if(T&&T.nodeType===1&&T.hasAttribute(Fn))return T}})(s),g=f!==void 0?f.nextSibling:null;u.setAttribute(Fn,"active"),u.setAttribute("data-styled-version","5.3.11");var h=Tm();return h&&u.setAttribute("nonce",h),s.insertBefore(u,g),u},jm=(function(){function n(s){var u=this.element=Cf(s);u.appendChild(document.createTextNode("")),this.sheet=(function(f){if(f.sheet)return f.sheet;for(var g=document.styleSheets,h=0,v=g.length;h<v;h++){var x=g[h];if(x.ownerNode===f)return x}sn(17)})(u),this.length=0}var l=n.prototype;return l.insertRule=function(s,u){try{return this.sheet.insertRule(u,s),this.length++,!0}catch{return!1}},l.deleteRule=function(s){this.sheet.deleteRule(s),this.length--},l.getRule=function(s){var u=this.sheet.cssRules[s];return u!==void 0&&typeof u.cssText=="string"?u.cssText:""},n})(),Pm=(function(){function n(s){var u=this.element=Cf(s);this.nodes=u.childNodes,this.length=0}var l=n.prototype;return l.insertRule=function(s,u){if(s<=this.length&&s>=0){var f=document.createTextNode(u),g=this.nodes[s];return this.element.insertBefore(f,g||null),this.length++,!0}return!1},l.deleteRule=function(s){this.element.removeChild(this.nodes[s]),this.length--},l.getRule=function(s){return s<this.length?this.nodes[s].textContent:""},n})(),Rm=(function(){function n(s){this.rules=[],this.length=0}var l=n.prototype;return l.insertRule=function(s,u){return s<=this.length&&(this.rules.splice(s,0,u),this.length++,!0)},l.deleteRule=function(s){this.rules.splice(s,1),this.length--},l.getRule=function(s){return s<this.length?this.rules[s]:""},n})(),Qd=Ds,$m={isServer:!Ds,useCSSOMInjection:!xm},tl=(function(){function n(s,u,f){s===void 0&&(s=Br),u===void 0&&(u={}),this.options=Zt({},$m,{},s),this.gs=u,this.names=new Map(f),this.server=!!s.isServer,!this.server&&Ds&&Qd&&(Qd=!1,(function(g){for(var h=document.querySelectorAll(Cm),v=0,x=h.length;v<x;v++){var E=h[v];E&&E.getAttribute(Fn)!=="active"&&(_m(g,E),E.parentNode&&E.parentNode.removeChild(E))}})(this))}n.registerId=function(s){return Vi(s)};var l=n.prototype;return l.reconstructWithOptions=function(s,u){return u===void 0&&(u=!0),new n(Zt({},this.options,{},s),this.gs,u&&this.names||void 0)},l.allocateGSInstance=function(s){return this.gs[s]=(this.gs[s]||0)+1},l.getTag=function(){return this.tag||(this.tag=(f=(u=this.options).isServer,g=u.useCSSOMInjection,h=u.target,s=f?new Rm(h):g?new jm(h):new Pm(h),new bm(s)));var s,u,f,g,h},l.hasNameForId=function(s,u){return this.names.has(s)&&this.names.get(s).has(u)},l.registerName=function(s,u){if(Vi(s),this.names.has(s))this.names.get(s).add(u);else{var f=new Set;f.add(u),this.names.set(s,f)}},l.insertRules=function(s,u,f){this.registerName(s,u),this.getTag().insertRules(Vi(s),f)},l.clearNames=function(s){this.names.has(s)&&this.names.get(s).clear()},l.clearRules=function(s){this.getTag().clearGroup(Vi(s)),this.clearNames(s)},l.clearTag=function(){this.tag=void 0},l.toString=function(){return(function(s){for(var u=s.getTag(),f=u.length,g="",h=0;h<f;h++){var v=km(h);if(v!==void 0){var x=s.names.get(v),E=u.getGroup(h);if(x&&E&&x.size){var T=Fn+".g"+h+'[id="'+v+'"]',A="";x!==void 0&&x.forEach((function(P){P.length>0&&(A+=P+",")})),g+=""+E+T+'{content:"'+A+`"}/*!sc*/
`}}}return g})(this)},n})(),Dm=/(a)(d)/gi,Gd=function(n){return String.fromCharCode(n+(n>25?39:97))};function ks(n){var l,s="";for(l=Math.abs(n);l>52;l=l/52|0)s=Gd(l%52)+s;return(Gd(l%52)+s).replace(Dm,"$1-$2")}var Mn=function(n,l){for(var s=l.length;s;)n=33*n^l.charCodeAt(--s);return n},Af=function(n){return Mn(5381,n)};function Ef(n){for(var l=0;l<n.length;l+=1){var s=n[l];if(Bn(s)&&!$s(s))return!1}return!0}var zm=Af("5.3.11"),Im=(function(){function n(l,s,u){this.rules=l,this.staticRulesId="",this.isStatic=(u===void 0||u.isStatic)&&Ef(l),this.componentId=s,this.baseHash=Mn(zm,s),this.baseStyle=u,tl.registerId(s)}return n.prototype.generateAndInjectStyles=function(l,s,u){var f=this.componentId,g=[];if(this.baseStyle&&g.push(this.baseStyle.generateAndInjectStyles(l,s,u)),this.isStatic&&!u.hash)if(this.staticRulesId&&s.hasNameForId(f,this.staticRulesId))g.push(this.staticRulesId);else{var h=cn(this.rules,l,s,u).join(""),v=ks(Mn(this.baseHash,h)>>>0);if(!s.hasNameForId(f,v)){var x=u(h,"."+v,void 0,f);s.insertRules(f,v,x)}g.push(v),this.staticRulesId=v}else{for(var E=this.rules.length,T=Mn(this.baseHash,u.hash),A="",P=0;P<E;P++){var U=this.rules[P];if(typeof U=="string")A+=U;else if(U){var j=cn(U,l,s,u),y=Array.isArray(j)?j.join(""):j;T=Mn(T,y+P),A+=y}}if(A){var $=ks(T>>>0);if(!s.hasNameForId(f,$)){var Y=u(A,"."+$,void 0,f);s.insertRules(f,$,Y)}g.push($)}}return g.join(" ")},n})(),Lm=/^\s*\/\/.*$/gm,Mm=[":","[",".","#"];function Om(n){var l,s,u,f,g=Br,h=g.options,v=h===void 0?Br:h,x=g.plugins,E=x===void 0?Ji:x,T=new um(v),A=[],P=(function(y){function $(Y){if(Y)try{y(Y+"}")}catch{}}return function(Y,H,V,J,L,M,ne,de,Ae,$e){switch(Y){case 1:if(Ae===0&&H.charCodeAt(0)===64)return y(H+";"),"";break;case 2:if(de===0)return H+"/*|*/";break;case 3:switch(de){case 102:case 112:return y(V[0]+H),"";default:return H+($e===0?"/*|*/":"")}case-2:H.split("/*|*/}").forEach($)}}})((function(y){A.push(y)})),U=function(y,$,Y){return $===0&&Mm.indexOf(Y[s.length])!==-1||Y.match(f)?y:"."+l};function j(y,$,Y,H){H===void 0&&(H="&");var V=y.replace(Lm,""),J=$&&Y?Y+" "+$+" { "+V+" }":V;return l=H,s=$,u=new RegExp("\\"+s+"\\b","g"),f=new RegExp("(\\"+s+"\\b){2,}"),T(Y||!$?"":$,J)}return T.use([].concat(E,[function(y,$,Y){y===2&&Y.length&&Y[0].lastIndexOf(s)>0&&(Y[0]=Y[0].replace(u,U))},P,function(y){if(y===-2){var $=A;return A=[],$}}])),j.hash=E.length?E.reduce((function(y,$){return $.name||sn(15),Mn(y,$.name)}),5381).toString():"",j}var _f=B.createContext();_f.Consumer;var Tf=B.createContext(),Nm=(Tf.Consumer,new tl),Ss=Om();function jf(){return m.useContext(_f)||Nm}function Pf(){return m.useContext(Tf)||Ss}var Bm=(function(){function n(l,s){var u=this;this.inject=function(f,g){g===void 0&&(g=Ss);var h=u.name+g.hash;f.hasNameForId(u.id,h)||f.insertRules(u.id,h,g(u.rules,h,"@keyframes"))},this.toString=function(){return sn(12,String(u.name))},this.name=l,this.id="sc-keyframes-"+l,this.rules=s}return n.prototype.getName=function(l){return l===void 0&&(l=Ss),this.name+l.hash},n})(),Fm=/([A-Z])/,Wm=/([A-Z])/g,Hm=/^ms-/,Um=function(n){return"-"+n.toLowerCase()};function Kd(n){return Fm.test(n)?n.replace(Wm,Um).replace(Hm,"-ms-"):n}var Yd=function(n){return n==null||n===!1||n===""};function cn(n,l,s,u){if(Array.isArray(n)){for(var f,g=[],h=0,v=n.length;h<v;h+=1)(f=cn(n[h],l,s,u))!==""&&(Array.isArray(f)?g.push.apply(g,f):g.push(f));return g}if(Yd(n))return"";if($s(n))return"."+n.styledComponentId;if(Bn(n)){if(typeof(E=n)!="function"||E.prototype&&E.prototype.isReactComponent||!l)return n;var x=n(l);return cn(x,l,s,u)}var E;return n instanceof Bm?s?(n.inject(s,u),n.getName(u)):n:bs(n)?(function T(A,P){var U,j,y=[];for(var $ in A)A.hasOwnProperty($)&&!Yd(A[$])&&(Array.isArray(A[$])&&A[$].isCss||Bn(A[$])?y.push(Kd($)+":",A[$],";"):bs(A[$])?y.push.apply(y,T(A[$],$)):y.push(Kd($)+": "+(U=$,(j=A[$])==null||typeof j=="boolean"||j===""?"":typeof j!="number"||j===0||U in dm||U.startsWith("--")?String(j).trim():j+"px")+";"));return P?[P+" {"].concat(y,["}"]):y})(n):n.toString()}var qd=function(n){return Array.isArray(n)&&(n.isCss=!0),n};function me(n){for(var l=arguments.length,s=new Array(l>1?l-1:0),u=1;u<l;u++)s[u-1]=arguments[u];return Bn(n)||bs(n)?qd(cn(Ud(Ji,[n].concat(s)))):s.length===0&&n.length===1&&typeof n[0]=="string"?n:qd(cn(Ud(n,s)))}var Rf=function(n,l,s){return s===void 0&&(s=Br),n.theme!==s.theme&&n.theme||l||s.theme},Vm=/[!"#$%&'()*+,./:;<=>?@[\\\]^`{|}~-]+/g,Qm=/(^-|-$)/g;function ns(n){return n.replace(Vm,"-").replace(Qm,"")}var $f=function(n){return ks(Af(n)>>>0)};function Qi(n){return typeof n=="string"&&!0}var Cs=function(n){return typeof n=="function"||typeof n=="object"&&n!==null&&!Array.isArray(n)},Gm=function(n){return n!=="__proto__"&&n!=="constructor"&&n!=="prototype"};function Km(n,l,s){var u=n[s];Cs(l)&&Cs(u)?Df(u,l):n[s]=l}function Df(n){for(var l=arguments.length,s=new Array(l>1?l-1:0),u=1;u<l;u++)s[u-1]=arguments[u];for(var f=0,g=s;f<g.length;f++){var h=g[f];if(Cs(h))for(var v in h)Gm(v)&&Km(n,h[v],v)}return n}var Ro=B.createContext();Ro.Consumer;function Ym(n){var l=m.useContext(Ro),s=m.useMemo((function(){return(function(u,f){if(!u)return sn(14);if(Bn(u)){var g=u(f);return g}return Array.isArray(u)||typeof u!="object"?sn(8):f?Zt({},f,{},u):u})(n.theme,l)}),[n.theme,l]);return n.children?B.createElement(Ro.Provider,{value:s},n.children):null}var os={};function zf(n,l,s){var u=$s(n),f=!Qi(n),g=l.attrs,h=g===void 0?Ji:g,v=l.componentId,x=v===void 0?(function(H,V){var J=typeof H!="string"?"sc":ns(H);os[J]=(os[J]||0)+1;var L=J+"-"+$f("5.3.11"+J+os[J]);return V?V+"-"+L:L})(l.displayName,l.parentComponentId):v,E=l.displayName,T=E===void 0?(function(H){return Qi(H)?"styled."+H:"Styled("+Vd(H)+")"})(n):E,A=l.displayName&&l.componentId?ns(l.displayName)+"-"+l.componentId:l.componentId||x,P=u&&n.attrs?Array.prototype.concat(n.attrs,h).filter(Boolean):h,U=l.shouldForwardProp;u&&n.shouldForwardProp&&(U=l.shouldForwardProp?function(H,V,J){return n.shouldForwardProp(H,V,J)&&l.shouldForwardProp(H,V,J)}:n.shouldForwardProp);var j,y=new Im(s,A,u?n.componentStyle:void 0),$=y.isStatic&&h.length===0,Y=function(H,V){return(function(J,L,M,ne){var de=J.attrs,Ae=J.componentStyle,$e=J.defaultProps,De=J.foldedComponentIds,W=J.shouldForwardProp,ce=J.styledComponentId,xe=J.target,K=(function(q,b,oe){q===void 0&&(q=Br);var Q=Zt({},b,{theme:q}),_={};return oe.forEach((function(le){var fe,ae,je,ze=le;for(fe in Bn(ze)&&(ze=ze(Q)),ze)Q[fe]=_[fe]=fe==="className"?(ae=_[fe],je=ze[fe],ae&&je?ae+" "+je:ae||je):ze[fe]})),[Q,_]})(Rf(L,m.useContext(Ro),$e)||Br,L,de),ue=K[0],ge=K[1],F=(function(q,b,oe,Q){var _=jf(),le=Pf(),fe=b?q.generateAndInjectStyles(Br,_,le):q.generateAndInjectStyles(oe,_,le);return fe})(Ae,ne,ue),se=M,G=ge.$as||L.$as||ge.as||L.as||xe,S=Qi(G),C=ge!==L?Zt({},L,{},ge):L,I={};for(var D in C)D[0]!=="$"&&D!=="as"&&(D==="forwardedAs"?I.as=C[D]:(W?W(D,Bd,G):!S||Bd(D))&&(I[D]=C[D]));return L.style&&ge.style!==L.style&&(I.style=Zt({},L.style,{},ge.style)),I.className=Array.prototype.concat(De,ce,F!==ce?F:null,L.className,ge.className).filter(Boolean).join(" "),I.ref=se,m.createElement(G,I)})(j,H,V,$)};return Y.displayName=T,(j=B.forwardRef(Y)).attrs=P,j.componentStyle=y,j.displayName=T,j.shouldForwardProp=U,j.foldedComponentIds=u?Array.prototype.concat(n.foldedComponentIds,n.styledComponentId):Ji,j.styledComponentId=A,j.target=u?n.target:n,j.withComponent=function(H){var V=l.componentId,J=(function(M,ne){if(M==null)return{};var de,Ae,$e={},De=Object.keys(M);for(Ae=0;Ae<De.length;Ae++)de=De[Ae],ne.indexOf(de)>=0||($e[de]=M[de]);return $e})(l,["componentId"]),L=V&&V+"-"+(Qi(H)?H:ns(Vd(H)));return zf(H,Zt({},J,{attrs:P,componentId:L}),s)},Object.defineProperty(j,"defaultProps",{get:function(){return this._foldedDefaultProps},set:function(H){this._foldedDefaultProps=u?Df({},n.defaultProps,H):H}}),Object.defineProperty(j,"toString",{value:function(){return"."+j.styledComponentId}}),f&&ym(j,n,{attrs:!0,componentStyle:!0,displayName:!0,foldedComponentIds:!0,shouldForwardProp:!0,styledComponentId:!0,target:!0,withComponent:!0}),j}var ee=function(n){return(function l(s,u,f){if(f===void 0&&(f=Br),!Sf.isValidElementType(u))return sn(1,String(u));var g=function(){return s(u,f,me.apply(void 0,arguments))};return g.withConfig=function(h){return l(s,u,Zt({},f,{},h))},g.attrs=function(h){return l(s,u,Zt({},f,{attrs:Array.prototype.concat(f.attrs,h).filter(Boolean)}))},g})(zf,n)};["a","abbr","address","area","article","aside","audio","b","base","bdi","bdo","big","blockquote","body","br","button","canvas","caption","cite","code","col","colgroup","data","datalist","dd","del","details","dfn","dialog","div","dl","dt","em","embed","fieldset","figcaption","figure","footer","form","h1","h2","h3","h4","h5","h6","head","header","hgroup","hr","html","i","iframe","img","input","ins","kbd","keygen","label","legend","li","link","main","map","mark","marquee","menu","menuitem","meta","meter","nav","noscript","object","ol","optgroup","option","output","p","param","picture","pre","progress","q","rp","rt","ruby","s","samp","script","section","select","small","source","span","strong","style","sub","summary","sup","table","tbody","td","textarea","tfoot","th","thead","time","title","tr","track","u","ul","var","video","wbr","circle","clipPath","defs","ellipse","foreignObject","g","image","line","linearGradient","marker","mask","path","pattern","polygon","polyline","radialGradient","rect","stop","svg","text","textPath","tspan"].forEach((function(n){ee[n]=ee(n)}));var qm=(function(){function n(s,u){this.rules=s,this.componentId=u,this.isStatic=Ef(s),tl.registerId(this.componentId+1)}var l=n.prototype;return l.createStyles=function(s,u,f,g){var h=g(cn(this.rules,u,f,g).join(""),""),v=this.componentId+s;f.insertRules(v,v,h)},l.removeStyles=function(s,u){u.clearRules(this.componentId+s)},l.renderStyles=function(s,u,f,g){s>2&&tl.registerId(this.componentId+s),this.removeStyles(s,f),this.createStyles(s,u,f,g)},n})();function Xm(n){for(var l=arguments.length,s=new Array(l>1?l-1:0),u=1;u<l;u++)s[u-1]=arguments[u];var f=me.apply(void 0,[n].concat(s)),g="sc-global-"+$f(JSON.stringify(f)),h=new qm(f,g);function v(E){var T=jf(),A=Pf(),P=m.useContext(Ro),U=m.useRef(T.allocateGSInstance(g)).current;return T.server&&x(U,E,T,P,A),m.useLayoutEffect((function(){if(!T.server)return x(U,E,T,P,A),function(){return h.removeStyles(U,T)}}),[U,E,T,P,A]),null}function x(E,T,A,P,U){if(h.isStatic)h.renderStyles(E,wm,A,U);else{var j=Zt({},T,{theme:Rf(T,P,v.defaultProps)});h.renderStyles(E,j,A,U)}}return B.memo(v)}const zs="4px 4px 10px 0 rgba(0, 0, 0, 0.35)",Is="inset 2px 2px 3px rgba(0,0,0,0.2)",Jt=()=>me`
  -webkit-text-fill-color: ${({theme:n})=>n.materialTextDisabled};
  color: ${({theme:n})=>n.materialTextDisabled};
  text-shadow: 1px 1px ${({theme:n})=>n.materialTextDisabledShadow};
  /* filter: grayscale(100%); */
`,er=({background:n="material",color:l="materialText"}={})=>me`
  box-sizing: border-box;
  display: inline-block;
  background: ${({theme:s})=>s[n]};
  color: ${({theme:s})=>s[l]};
`,Do=({mainColor:n="black",secondaryColor:l="transparent",pixelSize:s=2})=>me`
  background-image: ${[`linear-gradient(
      45deg,
      ${n} 25%,
      transparent 25%,
      transparent 75%,
      ${n} 75%
    )`,`linear-gradient(
      45deg,
      ${n} 25%,
      transparent 25%,
      transparent 75%,
      ${n} 75%
    )`].join(",")};
  background-color: ${l};
  background-size: ${`${s*2}px ${s*2}px`};
  background-position: 0 0, ${`${s}px ${s}px`};
`,un=()=>me`
  position: relative;
  box-sizing: border-box;
  display: inline-block;
  color: ${({theme:n})=>n.materialText};
  background: ${({$disabled:n,theme:l})=>n?l.flatLight:l.canvas};
  border: 2px solid ${({theme:n})=>n.canvas};
  outline: 2px solid ${({theme:n})=>n.flatDark};
  outline-offset: -4px;
`,Ln={button:{topLeftOuter:"borderLightest",topLeftInner:"borderLight",bottomRightInner:"borderDark",bottomRightOuter:"borderDarkest"},buttonPressed:{topLeftOuter:"borderDarkest",topLeftInner:"borderDark",bottomRightInner:"borderLight",bottomRightOuter:"borderLightest"},buttonThin:{topLeftOuter:"borderLightest",topLeftInner:null,bottomRightInner:null,bottomRightOuter:"borderDark"},buttonThinPressed:{topLeftOuter:"borderDark",topLeftInner:null,bottomRightInner:null,bottomRightOuter:"borderLightest"},field:{topLeftOuter:"borderDark",topLeftInner:"borderDarkest",bottomRightInner:"borderLight",bottomRightOuter:"borderLightest"},grouping:{topLeftOuter:"borderDark",topLeftInner:"borderLightest",bottomRightInner:"borderDark",bottomRightOuter:"borderLightest"},status:{topLeftOuter:"borderDark",topLeftInner:null,bottomRightInner:null,bottomRightOuter:"borderLightest"},window:{topLeftOuter:"borderLight",topLeftInner:"borderLightest",bottomRightInner:"borderDark",bottomRightOuter:"borderDarkest"}},Zm=({theme:n,topLeftInner:l,bottomRightInner:s,hasShadow:u=!1,hasInsetShadow:f=!1})=>[u?zs:!1,f?Is:!1,l!==null?`inset 1px 1px 0px 1px ${n[l]}`:!1,s!==null?`inset -1px -1px 0 1px ${n[s]}`:!1].filter(Boolean).join(", "),nt=({invert:n=!1,style:l="button"}={})=>{const s={topLeftOuter:n?"bottomRightOuter":"topLeftOuter",topLeftInner:n?"bottomRightInner":"topLeftInner",bottomRightInner:n?"topLeftInner":"bottomRightInner",bottomRightOuter:n?"topLeftOuter":"bottomRightOuter"};return me`
    border-style: solid;
    border-width: 2px;
    border-left-color: ${({theme:u})=>u[Ln[l][s.topLeftOuter]]};
    border-top-color: ${({theme:u})=>u[Ln[l][s.topLeftOuter]]};
    border-right-color: ${({theme:u})=>u[Ln[l][s.bottomRightOuter]]};
    border-bottom-color: ${({theme:u})=>u[Ln[l][s.bottomRightOuter]]};
    box-shadow: ${({theme:u,shadow:f})=>Zm({theme:u,topLeftInner:Ln[l][s.topLeftInner],bottomRightInner:Ln[l][s.bottomRightInner],hasShadow:f})};
  `},Wn=()=>me`
  outline: 2px dotted ${({theme:n})=>n.materialText};
`,Jm=n=>Buffer.from(n).toString("base64"),eg=typeof btoa<"u"?btoa:Jm,Gi=(n,l=0)=>{const s=`<svg height="26" width="26" viewBox="0 0 26 26" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <g transform="rotate(${l} 13 13)">
      <polygon fill="${n}" points="6,10 20,10 13,17"/>
    </g>
  </svg>`;return`url(data:image/svg+xml;base64,${eg(s)})`},Ls=(n="default")=>me`
  ::-webkit-scrollbar {
    width: 26px;
    height: 26px;
  }
  ::-webkit-scrollbar-track {
    ${({theme:l})=>Do({mainColor:n==="flat"?l.flatLight:l.material,secondaryColor:n==="flat"?l.canvas:l.borderLightest})}
  }
  ::-webkit-scrollbar-thumb {
    ${er()}
    ${n==="flat"?un():nt({style:"window"})}
      outline-offset: -2px;
  }

  ::-webkit-scrollbar-corner {
    background-color: ${({theme:l})=>l.material};
  }
  ::-webkit-scrollbar-button {
    ${er()}
    ${n==="flat"?un():nt({style:"window"})}
      display: block;
    outline-offset: -2px;
    height: 26px;
    width: 26px;
    background-repeat: no-repeat;
    background-size: 100%;
    background-position: 0 0;
  }
  ::-webkit-scrollbar-button:active,
  ::-webkit-scrollbar-button:active {
    background-position: 0 1px;
    ${n==="default"?nt({style:"window",invert:!0}):""}
  }

  ::-webkit-scrollbar-button:horizontal:increment:start,
  ::-webkit-scrollbar-button:horizontal:decrement:end,
  ::-webkit-scrollbar-button:vertical:increment:start,
  ::-webkit-scrollbar-button:vertical:decrement:end {
    display: none;
  }

  ::-webkit-scrollbar-button:horizontal:decrement {
    background-image: ${({theme:l})=>Gi(l.materialText,90)};
  }

  ::-webkit-scrollbar-button:horizontal:increment {
    background-image: ${({theme:l})=>Gi(l.materialText,270)};
  }

  ::-webkit-scrollbar-button:vertical:decrement {
    background-image: ${({theme:l})=>Gi(l.materialText,180)};
  }

  ::-webkit-scrollbar-button:vertical:increment {
    background-image: ${({theme:l})=>Gi(l.materialText,0)};
  }
`,tg=ee.a`
  color: ${({theme:n})=>n.anchor};
  font-size: inherit;
  text-decoration: ${({underline:n})=>n?"underline":"none"};
  &:visited {
    color: ${({theme:n})=>n.anchorVisited};
  }
`,rg=m.forwardRef(({children:n,underline:l=!0,...s},u)=>B.createElement(tg,{ref:u,underline:l,...s},n));rg.displayName="Anchor";const ng=ee.header`
  ${nt()};
  ${er()};

  position: ${n=>{var l;return(l=n.position)!==null&&l!==void 0?l:n.fixed?"fixed":"absolute"}};
  top: 0;
  right: 0;
  left: auto;
  display: flex;
  flex-direction: column;
  width: 100%;
`,If=m.forwardRef(({children:n,fixed:l=!0,position:s="fixed",...u},f)=>B.createElement(ng,{fixed:l,position:l!==!1?s:void 0,ref:f,...u},n));If.displayName="AppBar";const dn=()=>{};function ln(n,l,s){return s!==null&&n>s?s:l!==null&&n<l?l:n}function og(n){if(Math.abs(n)<1){const s=n.toExponential().split("e-"),u=s[0].split(".")[1];return(u?u.length:0)+parseInt(s[1],10)}const l=n.toString().split(".")[1];return l?l.length:0}function Xd(n,l,s){const u=Math.round((n-s)/l)*l+s;return Number(u.toFixed(og(l)))}function Fr(n){return typeof n=="number"?`${n}px`:n}const ig=ee.div`
  display: inline-block;
  box-sizing: border-box;
  object-fit: contain;
  ${({size:n})=>`
    height: ${n};
    width: ${n};
    `}
  border-radius: ${({square:n})=>n?0:"50%"};
  overflow: hidden;
  ${({noBorder:n,theme:l})=>!n&&`
    border-top: 2px solid ${l.borderDark};
    border-left: 2px solid ${l.borderDark};
    border-bottom: 2px solid ${l.borderLightest};
    border-right: 2px solid ${l.borderLightest};
    background: ${l.material};
  `}
  ${({src:n})=>!n&&`
    display: flex;
    align-items: center;
    justify-content: space-around;
    font-weight: bold;
    font-size: 1rem;
  `}
`,lg=ee.img`
  display: block;
  object-fit: contain;
  width: 100%;
  height: 100%;
`,ag=m.forwardRef(({alt:n="",children:l,noBorder:s=!1,size:u=35,square:f=!1,src:g,...h},v)=>B.createElement(ig,{noBorder:s,ref:v,size:Fr(u),square:f,src:g,...h},g?B.createElement(lg,{src:g,alt:n}):l));ag.displayName="Avatar";const gt={sm:"28px",md:"36px",lg:"44px"},sg=me`
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: ${({size:n="md"})=>gt[n]};
  width: ${({fullWidth:n,size:l="md",square:s})=>n?"100%":s?gt[l]:"auto"};
  padding: ${({square:n})=>n?0:"0 10px"};
  font-size: 1rem;
  user-select: none;
  &:active {
    padding-top: ${({disabled:n})=>!n&&"2px"};
  }
  padding-top: ${({active:n,disabled:l})=>n&&!l&&"2px"};
  &:after {
    content: '';
    position: absolute;
    display: block;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
  }
  &:not(:disabled) {
    cursor: pointer;
  }
  font-family: inherit;
`,nl=ee.button`
  ${({active:n,disabled:l,primary:s,theme:u,variant:f})=>f==="flat"?me`
          ${un()}
          ${s?`
          border: 2px solid ${u.checkmark};
            outline: 2px solid ${u.flatDark};
            outline-offset: -4px;
          `:`
          border: 2px solid ${u.flatDark};
            outline: 2px solid transparent;
            outline-offset: -4px;
          `}
          &:focus:after, &:active:after {
            ${!n&&!l&&Wn}
            outline-offset: -4px;
          }
        `:f==="menu"||f==="thin"?me`
          ${er()};
          border: 2px solid transparent;
          &:hover,
          &:focus {
            ${!l&&!n&&nt({style:"buttonThin"})}
          }
          &:active {
            ${!l&&nt({style:"buttonThinPressed"})}
          }
          ${n&&nt({style:"buttonThinPressed"})}
          ${l&&Jt()}
        `:me`
          ${er()};
          border: none;
          ${l&&Jt()}
          ${n?Do({mainColor:u.material,secondaryColor:u.borderLightest}):""}
          &:before {
            box-sizing: border-box;
            content: '';
            position: absolute;
            ${s?me`
                  left: 2px;
                  top: 2px;
                  width: calc(100% - 4px);
                  height: calc(100% - 4px);
                  outline: 2px solid ${u.borderDarkest};
                `:me`
                  left: 0;
                  top: 0;
                  width: 100%;
                  height: 100%;
                `}

            ${nt(n?{style:f==="raised"?"window":"button",invert:!0}:{style:f==="raised"?"window":"button",invert:!1})}
          }
          &:active:before {
            ${!l&&nt({style:f==="raised"?"window":"button",invert:!0})}
          }
          &:focus:after,
          &:active:after {
            ${!n&&!l&&Wn}
            outline-offset: -8px;
          }
          &:active:focus:after,
          &:active:after {
            top: ${n?"0":"1px"};
          }
        `}
  ${sg}
`,he=m.forwardRef(({onClick:n,disabled:l=!1,children:s,type:u="button",fullWidth:f=!1,size:g="md",square:h=!1,active:v=!1,onTouchStart:x=dn,primary:E=!1,variant:T="default",...A},P)=>B.createElement(nl,{active:v,disabled:l,$disabled:l,fullWidth:f,onClick:l?void 0:n,onTouchStart:x,primary:E,ref:P,size:g,square:h,type:u,variant:T,...A},s));he.displayName="Button";function Wr({defaultValue:n,onChange:l,onChangePropName:s="onChange",readOnly:u,value:f,valuePropName:g="value"}){const h=f!==void 0,[v,x]=m.useState(n),E=m.useCallback(T=>{h||x(T)},[h]);if(h&&typeof l!="function"&&!u){const T=`Warning: You provided a \`${g}\` prop to a component without an \`${s}\` handler.${g==="value"?`This will render a read-only field. If the field should be mutable use \`defaultValue\`. Otherwise, set either \`${s}\` or \`readOnly\`.`:`This breaks the component state. You must provide an \`${s}\` function that updates \`${g}\`.`}`;console.warn(T)}return[h?f:v,E]}const As=ee.li`
  box-sizing: border-box;

  display: flex;
  align-items: center;
  position: relative;
  height: ${n=>gt[n.size]};
  width: ${n=>n.square?gt[n.size]:"auto"};
  padding: 0 8px;
  font-size: 1rem;
  white-space: nowrap;
  justify-content: ${n=>n.square?"space-around":"space-between"};
  text-align: center;
  line-height: ${n=>gt[n.size]};
  color: ${({theme:n})=>n.materialText};
  pointer-events: ${({$disabled:n})=>n?"none":"auto"};
  font-weight: ${({primary:n})=>n?"bold":"normal"};
  &:hover {
    ${({theme:n,$disabled:l})=>!l&&`
        color: ${n.materialTextInvert};
        background: ${n.hoverBackground};
      `}

    cursor: default;
  }
  ${n=>n.$disabled&&Jt()}
`,lr=m.forwardRef(({size:n="lg",disabled:l,square:s,children:u,onClick:f,primary:g,...h},v)=>B.createElement(As,{$disabled:l,size:n,square:s,onClick:l?void 0:f,primary:g,role:"menuitem",ref:v,"aria-disabled":l,...h},u));lr.displayName="MenuListItem";const Es=ee.ul.attrs(()=>({role:"menu"}))`
  box-sizing: border-box;
  width: ${n=>n.fullWidth?"100%":"auto"};
  padding: 4px;
  ${nt({style:"window"})}
  ${er()}
  ${n=>n.inline&&`
    display: inline-flex;
    align-items: center;
  `}
  list-style: none;
  position: relative;
`;Es.displayName="MenuList";const sr=20,rl=ee.input`
  position: absolute;
  left: 0;
  margin: 0;
  width: ${sr}px;
  height: ${sr}px;
  opacity: 0;
  z-index: -1;
`,Ms=ee.label`
  display: inline-flex;
  align-items: center;
  position: relative;
  margin: 8px 0;
  cursor: ${({$disabled:n})=>n?"auto":"pointer"};
  user-select: none;
  font-size: 1rem;
  color: ${({theme:n})=>n.materialText};
  ${n=>n.$disabled&&Jt()}

  ${As} & {
    margin: 0;
    height: 100%;
  }
  ${As}:hover & {
    ${({$disabled:n,theme:l})=>!n&&me`
        color: ${l.materialTextInvert};
      `};
  }
`,Os=ee.span`
  display: inline-block;
  line-height: 1;
  padding: 2px;
  ${rl}:focus ~ & {
    ${Wn}
  }
  ${rl}:not(:disabled) ~ &:active {
    ${Wn}
  }
`,cr=ee.div`
  position: relative;
  box-sizing: border-box;
  padding: 2px;
  font-size: 1rem;
  border-style: solid;
  border-width: 2px;
  border-left-color: ${({theme:n})=>n.borderDark};
  border-top-color: ${({theme:n})=>n.borderDark};
  border-right-color: ${({theme:n})=>n.borderLightest};
  border-bottom-color: ${({theme:n})=>n.borderLightest};
  line-height: 1.5;
  &:before {
    position: absolute;
    left: 0;
    top: 0;
    content: '';
    width: calc(100% - 4px);
    height: calc(100% - 4px);

    border-style: solid;
    border-width: 2px;
    border-left-color: ${({theme:n})=>n.borderDarkest};
    border-top-color: ${({theme:n})=>n.borderDarkest};
    border-right-color: ${({theme:n})=>n.borderLight};
    border-bottom-color: ${({theme:n})=>n.borderLight};

    pointer-events: none;
    ${n=>n.shadow&&`box-shadow:${Is};`}
  }
`,cg=ee.div`
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  padding: 4px;
  overflow: auto;
  ${Ls()}
`,wr=m.forwardRef(({children:n,shadow:l=!0,...s},u)=>B.createElement(cr,{ref:u,shadow:l,...s},B.createElement(cg,null,n)));wr.displayName="ScrollView";const Lf=me`
  width: ${sr}px;
  height: ${sr}px;
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin-right: 0.5rem;
`,ug=ee(cr)`
  ${Lf}
  width: ${sr}px;
  height: ${sr}px;
  background: ${({$disabled:n,theme:l})=>n?l.material:l.canvas};
  &:before {
    box-shadow: none;
  }
`,dg=ee.div`
  position: relative;
  box-sizing: border-box;
  display: inline-block;
  background: ${({$disabled:n,theme:l})=>n?l.flatLight:l.canvas};
  ${Lf}
  width: ${sr-4}px;
  height: ${sr-4}px;
  outline: none;
  border: 2px solid ${({theme:n})=>n.flatDark};
  background: ${({$disabled:n,theme:l})=>n?l.flatLight:l.canvas};
`,fg=ee.span.attrs(()=>({"data-testid":"checkmarkIcon"}))`
  display: inline-block;
  position: relative;
  width: 100%;
  height: 100%;
  &:after {
    content: '';
    display: block;
    position: absolute;
    left: 50%;
    top: calc(50% - 1px);
    width: 3px;
    height: 7px;

    border: solid
      ${({$disabled:n,theme:l})=>n?l.checkmarkDisabled:l.checkmark};
    border-width: 0 3px 3px 0;
    transform: translate(-50%, -50%) rotate(45deg);

    border-color: ${n=>n.$disabled?n.theme.checkmarkDisabled:n.theme.checkmark};
  }
`,pg=ee.span.attrs(()=>({"data-testid":"indeterminateIcon"}))`
  display: inline-block;
  position: relative;

  width: 100%;
  height: 100%;

  &:after {
    content: '';
    display: block;

    width: 100%;
    height: 100%;

    ${({$disabled:n,theme:l})=>Do({mainColor:n?l.checkmarkDisabled:l.checkmark})}
    background-position: 0px 0px, 2px 2px;
  }
`,hg={flat:dg,default:ug},ut=m.forwardRef(({checked:n,className:l="",defaultChecked:s=!1,disabled:u=!1,indeterminate:f=!1,label:g="",onChange:h=dn,style:v={},value:x,variant:E="default",...T},A)=>{var P;const[U,j]=Wr({defaultValue:s,onChange:h,readOnly:(P=T.readOnly)!==null&&P!==void 0?P:u,value:n}),y=m.useCallback(H=>{const V=H.target.checked;j(V),h(H)},[h,j]),$=hg[E];let Y=null;return f?Y=pg:U&&(Y=fg),B.createElement(Ms,{$disabled:u,className:l,style:v},B.createElement(rl,{disabled:u,onChange:u?void 0:y,readOnly:u,type:"checkbox",value:x,checked:U,"data-indeterminate":f,ref:A,...T}),B.createElement($,{$disabled:u,role:"presentation"},Y&&B.createElement(Y,{$disabled:u,variant:E})),g&&B.createElement(Os,null,g))});ut.displayName="Checkbox";const Po=ee.div`
  ${({orientation:n,theme:l,size:s="100%"})=>n==="vertical"?`
    height: ${Fr(s)};
    border-left: 2px solid ${l.borderDark};
    border-right: 2px solid ${l.borderLightest};
    margin: 0;
    `:`
    width: ${Fr(s)};
    border-bottom: 2px solid ${l.borderLightest};
    border-top: 2px solid ${l.borderDark};
    margin: 0;
    `}
`;Po.displayName="Separator";const mg=ee(nl)`
  padding-left: 8px;
`,gg=ee(Po)`
  height: 21px;
  position: relative;
  top: 0;
`,Mf=ee.input`
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  position: absolute;
  left: 0;
  top: 0;
  opacity: 0;
  z-index: 1;
  cursor: pointer;
  &:disabled {
    cursor: default;
  }
`,vg=ee.div`
  box-sizing: border-box;
  height: 19px;
  display: inline-block;
  width: 35px;
  margin-right: 5px;

  background: ${({color:n})=>n};

  ${({$disabled:n})=>n?me`
          border: 2px solid ${({theme:l})=>l.materialTextDisabled};
          filter: drop-shadow(
            1px 1px 0px ${({theme:l})=>l.materialTextDisabledShadow}
          );
        `:me`
          border: 2px solid ${({theme:l})=>l.materialText};
        `}
  ${Mf}:focus:not(:active) + &:after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    ${Wn}
    outline-offset: -8px;
  }
`,yg=ee.span`
  width: 0px;
  height: 0px;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  display: inline-block;
  margin-left: 6px;

  ${({$disabled:n})=>n?me`
          border-top: 6px solid ${({theme:l})=>l.materialTextDisabled};
          filter: drop-shadow(
            1px 1px 0px ${({theme:l})=>l.materialTextDisabledShadow}
          );
        `:me`
          border-top: 6px solid ${({theme:l})=>l.materialText};
        `}
  &:after {
    content: '';
    box-sizing: border-box;
    position: absolute;
    top: ${({variant:n})=>n==="flat"?"6px":"8px"};
    right: 8px;
    width: 16px;
    height: 19px;
  }
`,xg=m.forwardRef(({value:n,defaultValue:l,onChange:s=dn,disabled:u=!1,variant:f="default",...g},h)=>{var v;const[x,E]=Wr({defaultValue:l,onChange:s,readOnly:(v=g.readOnly)!==null&&v!==void 0?v:u,value:n}),T=A=>{const P=A.target.value;E(P),s(A)};return B.createElement(mg,{disabled:u,as:"div",variant:f,size:"md"},B.createElement(Mf,{onChange:T,readOnly:u,disabled:u,value:x??"#008080",type:"color",ref:h,...g}),B.createElement(vg,{$disabled:u,color:x??"#008080",role:"presentation"}),f==="default"&&B.createElement(gg,{orientation:"vertical"}),B.createElement(yg,{$disabled:u,variant:f}))});xg.displayName="ColorInput";const wg=ee.div`
  position: relative;
  --react95-digit-primary-color: #ff0102;
  --react95-digit-secondary-color: #740201;
  --react95-digit-bg-color: #000000;

  ${({pixelSize:n})=>me`
    width: ${11*n}px;
    height: ${21*n}px;
    margin: ${n}px;

    span,
    span:before,
    span:after {
      box-sizing: border-box;
      display: inline-block;
      position: absolute;
    }
    span.active,
    span.active:before,
    span.active:after {
      background: var(--react95-digit-primary-color);
    }
    span:not(.active),
    span:not(.active):before,
    span:not(.active):after {
      ${Do({mainColor:"var(--react95-digit-bg-color)",secondaryColor:"var(--react95-digit-secondary-color)",pixelSize:n})}
    }

    span.horizontal,
    span.horizontal:before,
    span.horizontal:after {
      height: ${n}px;
      border-left: ${n}px solid var(--react95-digit-bg-color);
      border-right: ${n}px solid var(--react95-digit-bg-color);
    }
    span.horizontal.active,
    span.horizontal.active:before,
    span.horizontal.active:after {
      height: ${n}px;
      border-left: ${n}px solid var(--react95-digit-primary-color);
      border-right: ${n}px solid var(--react95-digit-primary-color);
    }
    span.horizontal {
      left: ${n}px;
      width: ${9*n}px;
    }
    span.horizontal:before {
      content: '';
      width: 100%;
      top: ${n}px;
      left: ${0}px;
    }
    span.horizontal:after {
      content: '';
      width: calc(100% - ${n*2}px);
      top: ${2*n}px;
      left: ${n}px;
    }
    span.horizontal.top {
      top: 0;
    }
    span.horizontal.bottom {
      bottom: 0;
      transform: rotateX(180deg);
    }

    span.center,
    span.center:before,
    span.center:after {
      height: ${n}px;
      border-left: ${n}px solid var(--react95-digit-bg-color);
      border-right: ${n}px solid var(--react95-digit-bg-color);
    }
    span.center.active,
    span.center.active:before,
    span.center.active:after {
      border-left: ${n}px solid var(--react95-digit-primary-color);
      border-right: ${n}px solid var(--react95-digit-primary-color);
    }
    span.center {
      top: 50%;
      transform: translateY(-50%);
      left: ${n}px;
      width: ${9*n}px;
    }
    span.center:before,
    span.center:after {
      content: '';
      width: 100%;
    }
    span.center:before {
      top: ${n}px;
    }
    span.center:after {
      bottom: ${n}px;
    }

    span.vertical,
    span.vertical:before,
    span.vertical:after {
      width: ${n}px;
      border-top: ${n}px solid var(--react95-digit-bg-color);
      border-bottom: ${n}px solid var(--react95-digit-bg-color);
    }
    span.vertical {
      height: ${11*n}px;
    }
    span.vertical.left {
      left: 0;
    }
    span.vertical.right {
      right: 0;
      transform: rotateY(180deg);
    }
    span.vertical.top {
      top: 0px;
    }
    span.vertical.bottom {
      bottom: 0px;
    }
    span.vertical:before {
      content: '';
      height: 100%;
      top: ${0}px;
      left: ${n}px;
    }
    span.vertical:after {
      content: '';
      height: calc(100% - ${n*2}px);
      top: ${n}px;
      left: ${n*2}px;
    }
  `}
`,Zd=["horizontal top","center","horizontal bottom","vertical top left","vertical top right","vertical bottom left","vertical bottom right"],bg=[[1,0,1,1,1,1,1],[0,0,0,0,1,0,1],[1,1,1,0,1,1,0],[1,1,1,0,1,0,1],[0,1,0,1,1,0,1],[1,1,1,1,0,0,1],[1,1,1,1,0,1,1],[1,0,0,0,1,0,1],[1,1,1,1,1,1,1],[1,1,1,1,1,0,1]];function kg({digit:n=0,pixelSize:l=2,...s}){const u=bg[Number(n)].map((f,g)=>f?`${Zd[g]} active`:Zd[g]);return B.createElement(wg,{pixelSize:l,...s},u.map((f,g)=>B.createElement("span",{className:f,key:g})))}const Sg=ee.div`
  ${nt({style:"status"})}
  display: inline-flex;
  background: #000000;
`,Cg={sm:1,md:2,lg:3,xl:4},Ag=m.forwardRef(({value:n=0,minLength:l=3,size:s="md",...u},f)=>{const g=m.useMemo(()=>n.toString().padStart(l,"0").split(""),[l,n]);return B.createElement(Sg,{ref:f,...u},g.map((h,v)=>B.createElement(kg,{digit:h,pixelSize:Cg[s],key:v})))});Ag.displayName="Counter";const Of=me`
  display: flex;
  align-items: center;
  width: ${({fullWidth:n})=>n?"100%":"auto"};
  min-height: ${gt.md};
`,Eg=ee(cr).attrs({"data-testid":"variant-default"})`
  ${Of}
  background: ${({$disabled:n,theme:l})=>n?l.material:l.canvas};
`,_g=ee.div.attrs({"data-testid":"variant-flat"})`
  ${un()}
  ${Of}
  position: relative;
`,Nf=me`
  display: block;
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  outline: none;
  border: none;
  background: none;
  font-size: 1rem;
  min-height: 27px;
  font-family: inherit;
  color: ${({theme:n})=>n.canvasText};
  ${({disabled:n,variant:l})=>l!=="flat"&&n&&Jt()}
`,Tg=ee.input`
  ${Nf}
  padding: 0 8px;
`,jg=ee.textarea`
  ${Nf}
  padding: 8px;
  resize: none;
  ${({variant:n})=>Ls(n)}
`,Ke=m.forwardRef(({className:n,disabled:l=!1,fullWidth:s,onChange:u=dn,shadow:f=!0,style:g,variant:h="default",...v},x)=>{const E=h==="flat"?_g:Eg,T=m.useMemo(()=>{var A;return v.multiline?B.createElement(jg,{disabled:l,onChange:l?void 0:u,readOnly:l,ref:x,variant:h,...v}):B.createElement(Tg,{disabled:l,onChange:l?void 0:u,readOnly:l,ref:x,type:(A=v.type)!==null&&A!==void 0?A:"text",variant:h,...v})},[l,u,v,x,h]);return B.createElement(E,{className:n,fullWidth:s,$disabled:l,shadow:f,style:g},T)});Ke.displayName="TextInput";const Pg=ee.div`
  display: inline-flex;
  align-items: center;
`,_s=ee(he)`
  width: 30px;
  padding: 0;
  flex-shrink: 0;

  ${({variant:n})=>n==="flat"?me`
          height: calc(50% - 1px);
        `:me`
          height: 50%;
        `}
`,Rg=ee.div`
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap;
  justify-content: space-between;

  ${({variant:n})=>n==="flat"?me`
          height: calc(${gt.md} - 4px);
        `:me`
          height: ${gt.md};
          margin-left: 2px;
        `}
`,Jd=ee.span`
  width: 0px;
  height: 0px;
  display: inline-block;
  ${({invert:n})=>n?me`
          border-left: 4px solid transparent;
          border-right: 4px solid transparent;
          border-bottom: 4px solid ${({theme:l})=>l.materialText};
        `:me`
          border-left: 4px solid transparent;
          border-right: 4px solid transparent;
          border-top: 4px solid ${({theme:l})=>l.materialText};
        `}
  ${_s}:disabled & {
    filter: drop-shadow(
      1px 1px 0px ${({theme:n})=>n.materialTextDisabledShadow}
    );
    ${({invert:n})=>n?me`
            border-bottom-color: ${({theme:l})=>l.materialTextDisabled};
          `:me`
            border-top-color: ${({theme:l})=>l.materialTextDisabled};
          `}
  }
`,Bf=m.forwardRef(({className:n,defaultValue:l,disabled:s=!1,max:u,min:f,onChange:g,readOnly:h,step:v=1,style:x,value:E,variant:T="default",width:A,...P},U)=>{const[j,y]=Wr({defaultValue:l,onChange:g,readOnly:h,value:E}),$=m.useCallback(M=>{const ne=parseFloat(M.target.value);y(ne)},[y]),Y=m.useCallback(M=>{const ne=ln(parseFloat(((j??0)+M).toFixed(2)),f??null,u??null);y(ne),g==null||g(ne)},[u,f,g,y,j]),H=m.useCallback(()=>{j!==void 0&&(g==null||g(j))},[g,j]),V=m.useCallback(()=>{Y(v)},[Y,v]),J=m.useCallback(()=>{Y(-v)},[Y,v]),L=T==="flat"?"flat":"raised";return B.createElement(Pg,{className:n,style:{...x,width:A!==void 0?Fr(A):"auto"},...P},B.createElement(Ke,{value:j,variant:T,onChange:$,disabled:s,type:"number",readOnly:h,ref:U,fullWidth:!0,onBlur:H}),B.createElement(Rg,{variant:T},B.createElement(_s,{"data-testid":"increment",variant:L,disabled:s||h,onClick:V},B.createElement(Jd,{invert:!0})),B.createElement(_s,{"data-testid":"decrement",variant:L,disabled:s||h,onClick:J},B.createElement(Jd,null))))});Bf.displayName="NumberInput";function $g(){const n="0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";let l="";for(let s=0;s<10;s+=1)l+=n[Math.floor(Math.random()*n.length)];return l}const Ff=n=>m.useMemo(()=>$g(),[n]),Wf=me`
  box-sizing: border-box;
  padding-left: 4px;
  overflow: hidden;
  white-space: nowrap;
  user-select: none;
  line-height: 100%;
`,Hf=me`
  background: ${({theme:n})=>n.hoverBackground};
  color: ${({theme:n})=>n.canvasTextInvert};
`,Ns=ee.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  &:focus {
    outline: none;
  }
`,Dg=ee.div`
  ${Wf}
  padding-right: 8px;
  align-items: center;
  display: flex;
  height: calc(100% - 4px);
  width: calc(100% - 4px);
  margin: 0 2px;
  border: 2px solid transparent;
  ${Ns}:focus & {
    ${Hf}
    border: 2px dotted ${({theme:n})=>n.focusSecondary};
  }
`,Uf=me`
  height: ${gt.md};
  display: inline-block;
  color: ${({$disabled:n=!1,theme:l})=>n?Jt():l.canvasText};
  font-size: 1rem;
  cursor: ${({$disabled:n})=>n?"default":"pointer"};
`,zg=ee(cr)`
  ${Uf}
  background: ${({$disabled:n=!1,theme:l})=>n?l.material:l.canvas};
  &:focus {
    outline: 0;
  }
`,Ig=ee.div`
  ${un()}
  ${Uf}
  background: ${({$disabled:n=!1,theme:l})=>n?l.flatLight:l.canvas};
`,Lg=ee.select`
  -moz-appearance: none;
  -webkit-appearance: none;
  display: block;
  width: 100%;
  height: 100%;
  color: inherit;
  font-size: 1rem;
  border: 0;
  margin: 0;
  background: none;
  -webkit-tap-highlight-color: transparent;
  border-radius: 0;
  padding-right: 30px;
  ${Wf}
  cursor: pointer;
  &:disabled {
    ${Jt()};
    background: ${({theme:n})=>n.material};
    cursor: default;
  }
`,Vf=ee(nl).attrs(()=>({"aria-hidden":"true"}))`
  width: 30px;
  padding: 0;
  flex-shrink: 0;
  ${({variant:n="default"})=>n==="flat"?me`
          height: 100%;
          margin-right: 0;
        `:me`
          height: 100%;
        `}
  ${({native:n=!1,variant:l="default"})=>n&&(l==="flat"?`
      position: absolute;
      right: 0;
      height: 100%;
      `:`
    position: absolute;
    top: 2px;
    right: 2px;
    height: calc(100% - 4px);
    `)}
    pointer-events: ${({$disabled:n=!1,native:l=!1})=>n||l?"none":"auto"}
`,Mg=ee.span`
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 0;
  height: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  display: inline-block;
  border-top: 6px solid
    ${({$disabled:n=!1,theme:l})=>n?l.materialTextDisabled:l.materialText};
  ${({$disabled:n=!1,theme:l})=>n&&`
    filter: drop-shadow(1px 1px 0px ${l.materialTextDisabledShadow});
    border-top-color: ${l.materialTextDisabled};
    `}
  ${Vf}:active & {
    margin-top: 2px;
  }
`,Og=ee.ul`
  box-sizing: border-box;

  font-size: 1rem;
  position: absolute;
  transform: translateY(100%);
  left: 0;
  background: ${({theme:n})=>n.canvas};
  padding: 2px;
  border-top: none;
  cursor: default;
  z-index: 1;
  cursor: pointer;
  box-shadow: ${zs};
  ${({variant:n="default"})=>n==="flat"?me`
          bottom: 2px;
          width: 100%;
          border: 2px solid ${({theme:l})=>l.flatDark};
        `:me`
          bottom: -2px;
          width: calc(100% - 2px);
          border: 2px solid ${({theme:l})=>l.borderDarkest};
        `}
  ${({variant:n="default"})=>Ls(n)}
`,Ng=ee.li`
  box-sizing: border-box;

  width: 100%;
  padding-left: 8px;

  height: calc(${gt.md} - 4px);
  line-height: calc(${gt.md} - 4px);
  font-size: 1rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: ${({theme:n})=>n.canvasText};
  &:focus {
    outline: 0;
  }
  ${({active:n})=>n?Hf:""}
  user-select: none;
`,Bg=[],Qf=({className:n,defaultValue:l,disabled:s,native:u,onChange:f,options:g=Bg,readOnly:h,style:v,value:x,variant:E,width:T})=>{var A;const P=m.useMemo(()=>g.filter(Boolean),[g]),[U,j]=Wr({defaultValue:l??((A=P==null?void 0:P[0])===null||A===void 0?void 0:A.value),onChange:f,readOnly:h,value:x}),y=!(s||h),$=m.useMemo(()=>({className:n,style:{...v,width:T}}),[n,v,T]),Y=m.useMemo(()=>B.createElement(Vf,{as:"div","data-testid":"select-button",$disabled:s,native:u,tabIndex:-1,variant:E==="flat"?"flat":"raised"},B.createElement(Mg,{"data-testid":"select-icon",$disabled:s})),[s,u,E]),H=m.useMemo(()=>E==="flat"?Ig:zg,[E]);return m.useMemo(()=>({isEnabled:y,options:P,value:U,setValue:j,wrapperProps:$,DropdownButton:Y,Wrapper:H}),[Y,H,y,P,j,U,$])},Fg={ARROW_DOWN:"ArrowDown",ARROW_UP:"ArrowUp",END:"End",ENTER:"Enter",ESC:"Escape",HOME:"Home",SPACE:"Space",TAB:"Tab"},Wg=1e3,Hg=({onBlur:n,onChange:l,onClose:s,onFocus:u,onKeyDown:f,onMouseDown:g,onOpen:h,open:v,options:x,readOnly:E,value:T,selectRef:A,setValue:P,wrapperRef:U})=>{const j=m.useRef(null),y=m.useRef([]),$=m.useRef(0),Y=m.useRef(0),H=m.useRef(),V=m.useRef("search"),J=m.useRef(""),L=m.useRef(),[M,ne]=Wr({defaultValue:!1,onChange:h,onChangePropName:"onOpen",readOnly:E,value:v,valuePropName:"open"}),de=m.useMemo(()=>{const _=x.findIndex(le=>le.value===T);return $.current=ln(_,0,null),x[_]},[x,T]),[Ae,$e]=m.useState(x[0]),De=m.useCallback(_=>{const le=j.current,fe=y.current[_];if(!fe||!le){H.current=_;return}H.current=void 0;const ae=le.clientHeight,je=le.scrollTop,ze=le.scrollTop+ae,Ce=fe.offsetTop,Me=fe.offsetHeight,Re=fe.offsetTop+fe.offsetHeight;Ce<je&&le.scrollTo(0,Ce),Re>ze&&le.scrollTo(0,Ce-ae+Me),fe.focus({preventScroll:!0})},[j]),W=m.useCallback((_,{scroll:le}={})=>{var fe;const ae=x.length-1;let je;switch(_){case"first":{je=0;break}case"last":{je=ae;break}case"next":{je=ln(Y.current+1,0,ae);break}case"previous":{je=ln(Y.current-1,0,ae);break}case"selected":{je=ln((fe=$.current)!==null&&fe!==void 0?fe:0,0,ae);break}default:je=_}Y.current=je,$e(x[je]),le&&De(je)},[Y,x,De]),ce=m.useCallback(({fromEvent:_})=>{ne(!0),W("selected",{scroll:!0}),h==null||h({fromEvent:_})},[W,h,ne]),xe=m.useCallback(()=>{V.current="search",J.current="",clearTimeout(L.current)},[]),K=m.useCallback(({focusSelect:_,fromEvent:le})=>{var fe;s==null||s({fromEvent:le}),ne(!1),$e(x[0]),xe(),H.current=void 0,_&&((fe=A.current)===null||fe===void 0||fe.focus())},[xe,s,x,A,ne]),ue=m.useCallback(({fromEvent:_})=>{M?K({focusSelect:!1,fromEvent:_}):ce({fromEvent:_})},[K,ce,M]),ge=m.useCallback((_,{fromEvent:le})=>{$.current!==_&&($.current=_,P(x[_].value),l==null||l(x[_],{fromEvent:le}))},[l,x,P]),F=m.useCallback(({focusSelect:_,fromEvent:le})=>{ge(Y.current,{fromEvent:le}),K({focusSelect:_,fromEvent:le})},[K,ge]),se=m.useCallback((_,{fromEvent:le,select:fe})=>{var ae;switch(V.current==="cycleFirstLetter"&&_!==J.current&&(V.current="search"),_===J.current?V.current="cycleFirstLetter":J.current+=_,V.current){case"search":{let je=x.findIndex(ze=>{var Ce;return((Ce=ze.label)===null||Ce===void 0?void 0:Ce.toLocaleUpperCase().indexOf(J.current))===0});je<0&&(je=x.findIndex(ze=>{var Ce;return((Ce=ze.label)===null||Ce===void 0?void 0:Ce.toLocaleUpperCase().indexOf(_))===0}),J.current=_),je>=0&&(fe?ge(je,{fromEvent:le}):W(je,{scroll:!0}));break}case"cycleFirstLetter":{const je=fe?(ae=$.current)!==null&&ae!==void 0?ae:-1:Y.current;let ze=x.findIndex((Ce,Me)=>{var Re;return Me>je&&((Re=Ce.label)===null||Re===void 0?void 0:Re.toLocaleUpperCase().indexOf(_))===0});ze<0&&(ze=x.findIndex(Ce=>{var Me;return((Me=Ce.label)===null||Me===void 0?void 0:Me.toLocaleUpperCase().indexOf(_))===0})),ze>=0&&(fe?ge(ze,{fromEvent:le}):W(ze,{scroll:!0}));break}}clearTimeout(L.current),L.current=setTimeout(()=>{V.current==="search"&&(J.current="")},Wg)},[W,x,ge]),G=m.useCallback(_=>{var le;_.button===0&&(_.preventDefault(),(le=A.current)===null||le===void 0||le.focus(),ue({fromEvent:_}),g==null||g(_))},[g,A,ue]),S=m.useCallback(_=>{F({focusSelect:!0,fromEvent:_})},[F]),C=m.useCallback(_=>{const{altKey:le,code:fe,ctrlKey:ae,metaKey:je,shiftKey:ze}=_,{ARROW_DOWN:Ce,ARROW_UP:Me,END:Re,ENTER:Ne,ESC:ot,HOME:Fe,SPACE:N,TAB:te}=Fg,Ee=le||ae||je||ze;if(!(fe===te&&(le||ae||je)||fe!==te&&Ee))switch(fe){case Ce:{if(_.preventDefault(),!M){ce({fromEvent:_});return}W("next",{scroll:!0});break}case Me:{if(_.preventDefault(),!M){ce({fromEvent:_});return}W("previous",{scroll:!0});break}case Re:{if(_.preventDefault(),!M){ce({fromEvent:_});return}W("last",{scroll:!0});break}case Ne:{if(!M)return;_.preventDefault(),F({focusSelect:!0,fromEvent:_});break}case ot:{if(!M)return;_.preventDefault(),K({focusSelect:!0,fromEvent:_});break}case Fe:{if(_.preventDefault(),!M){ce({fromEvent:_});return}W("first",{scroll:!0});break}case N:{_.preventDefault(),M?F({focusSelect:!0,fromEvent:_}):ce({fromEvent:_});break}case te:{if(!M)return;ze||_.preventDefault(),F({focusSelect:!ze,fromEvent:_});break}default:!Ee&&fe.match(/^Key/)&&(_.preventDefault(),_.stopPropagation(),se(fe.replace(/^Key/,""),{select:!M,fromEvent:_}))}},[W,K,M,ce,se,F]),I=m.useCallback(_=>{C(_),f==null||f(_)},[C,f]),D=m.useCallback(_=>{W(_)},[W]),q=m.useCallback(_=>{M||(xe(),n==null||n(_))},[xe,n,M]),b=m.useCallback(_=>{xe(),u==null||u(_)},[xe,u]),oe=m.useCallback(_=>{j.current=_,H.current!==void 0&&De(H.current)},[De]),Q=m.useCallback((_,le)=>{y.current[le]=_,H.current===le&&De(H.current)},[De]);return m.useEffect(()=>{if(!M)return()=>{};const _=le=>{var fe;const ae=le.target;!((fe=U.current)===null||fe===void 0)&&fe.contains(ae)||(le.preventDefault(),K({focusSelect:!1,fromEvent:le}))};return document.addEventListener("mousedown",_),()=>{document.removeEventListener("mousedown",_)}},[K,M,U]),m.useMemo(()=>({activeOption:Ae,handleActivateOptionIndex:D,handleBlur:q,handleButtonKeyDown:I,handleDropdownKeyDown:C,handleFocus:b,handleMouseDown:G,handleOptionClick:S,handleSetDropdownRef:oe,handleSetOptionRef:Q,open:M,selectedOption:de}),[Ae,D,q,I,b,C,G,S,oe,Q,M,de])},Gf=m.forwardRef(({className:n,defaultValue:l,disabled:s,onChange:u,options:f,readOnly:g,style:h,value:v,variant:x,width:E,...T},A)=>{const{isEnabled:P,options:U,setValue:j,value:y,DropdownButton:$,Wrapper:Y}=Qf({defaultValue:l,disabled:s,native:!0,onChange:u,options:f,readOnly:g,value:v,variant:x}),H=m.useCallback(V=>{const J=U.find(L=>L.value===V.target.value);J&&(j(J.value),u==null||u(J,{fromEvent:V}))},[u,U,j]);return B.createElement(Y,{className:n,style:{...h,width:E}},B.createElement(Ns,null,B.createElement(Lg,{...T,disabled:s,onChange:P?H:dn,ref:A,value:y},U.map((V,J)=>{var L;return B.createElement("option",{key:`${V.value}-${J}`,value:V.value},(L=V.label)!==null&&L!==void 0?L:V.value)})),$))});Gf.displayName="SelectNative";function Ug({activateOptionIndex:n,active:l,index:s,onClick:u,option:f,selected:g,setRef:h}){const v=m.useCallback(()=>{n(s)},[n,s]),x=m.useCallback(T=>{h(T,s)},[s,h]),E=Ff();return B.createElement(Ng,{active:l,"aria-selected":g?"true":void 0,"data-value":f.value,id:E,onClick:u,onMouseEnter:v,ref:x,role:"option",tabIndex:0},f.label)}function Vg({"aria-label":n,"aria-labelledby":l,className:s,defaultValue:u,disabled:f=!1,formatDisplay:g,inputProps:h,labelId:v,menuMaxHeight:x,name:E,onBlur:T,onChange:A,onClose:P,onFocus:U,onKeyDown:j,onMouseDown:y,onOpen:$,open:Y,options:H,readOnly:V,shadow:J=!0,style:L,variant:M="default",value:ne,width:de="auto",...Ae},$e){const{isEnabled:De,options:W,setValue:ce,value:xe,wrapperProps:K,DropdownButton:ue,Wrapper:ge}=Qf({className:s,defaultValue:u,disabled:f,native:!1,onChange:A,options:H,style:L,readOnly:V,value:ne,variant:M,width:de}),F=m.useRef(null),se=m.useRef(null),G=m.useRef(null),{activeOption:S,handleActivateOptionIndex:C,handleBlur:I,handleButtonKeyDown:D,handleDropdownKeyDown:q,handleFocus:b,handleMouseDown:oe,handleOptionClick:Q,handleSetDropdownRef:_,handleSetOptionRef:le,open:fe,selectedOption:ae}=Hg({onBlur:T,onChange:A,onClose:P,onFocus:U,onKeyDown:j,onMouseDown:y,onOpen:$,open:Y,options:W,value:xe,selectRef:se,setValue:ce,wrapperRef:G});m.useImperativeHandle($e,()=>({focus:Ne=>{var ot;(ot=se.current)===null||ot===void 0||ot.focus(Ne)},node:F.current,value:String(xe)}),[xe]);const je=m.useMemo(()=>ae?typeof g=="function"?g(ae):ae.label:"",[g,ae]),ze=De?1:void 0,Ce=m.useMemo(()=>x?{overflow:"auto",maxHeight:x}:void 0,[x]),Me=Ff(),Re=m.useMemo(()=>W.map((Ne,ot)=>{const Fe=`${xe}-${ot}`,N=Ne===S,te=Ne===ae;return B.createElement(Ug,{activateOptionIndex:C,active:N,index:ot,key:Fe,onClick:Q,option:Ne,selected:te,setRef:le})}),[S,C,Q,le,W,ae,xe]);return B.createElement(ge,{...K,$disabled:f,ref:G,shadow:J,style:{...L,width:de}},B.createElement("input",{name:E,ref:F,type:"hidden",value:String(xe),...h}),B.createElement(Ns,{"aria-disabled":f,"aria-expanded":fe,"aria-haspopup":"listbox","aria-label":n,"aria-labelledby":l??v,"aria-owns":De&&fe?Me:void 0,onBlur:I,onFocus:b,onKeyDown:D,onMouseDown:De?oe:y,ref:se,role:"button",tabIndex:ze,...Ae},B.createElement(Dg,null,je),ue),De&&fe&&B.createElement(Og,{id:Me,onKeyDown:q,ref:_,role:"listbox",style:Ce,tabIndex:0,variant:M},Re))}const Nn=m.forwardRef(Vg);Nn.displayName="Select";const Qg=ee.div`
  position: relative;
  display: flex;
  align-items: center;
  padding: ${n=>n.noPadding?"0":"4px"};
`,$o=m.forwardRef(function({children:l,noPadding:s=!1,...u},f){return B.createElement(Qg,{noPadding:s,ref:f,...u},l)});$o.displayName="Toolbar";const Gg=ee.div`
  padding: 16px;
`,Bs=m.forwardRef(function({children:l,...s},u){return B.createElement(Gg,{ref:u,...s},l)});Bs.displayName="WindowContent";const Kg=ee.div`
  height: 33px;
  line-height: 33px;
  padding-left: 0.25rem;
  padding-right: 3px;
  font-weight: bold;
  border: 2px solid ${({theme:n})=>n.material};
  ${({active:n})=>n===!1?me`
          background: ${({theme:l})=>l.headerNotActiveBackground};
          color: ${({theme:l})=>l.headerNotActiveText};
        `:me`
          background: ${({theme:l})=>l.headerBackground};
          color: ${({theme:l})=>l.headerText};
        `}

  ${nl} {
    padding-left: 0;
    padding-right: 0;
    height: 27px;
    width: 31px;
  }
`,Fs=m.forwardRef(function({active:l=!0,children:s,...u},f){return B.createElement(Kg,{active:l,ref:f,...u},s)});Fs.displayName="WindowHeader";const Yg=ee.div`
  position: relative;
  padding: 4px;
  font-size: 1rem;
  ${nt({style:"window"})}
  ${er()}
`,qg=ee.span`
  ${({theme:n})=>me`
    display: inline-block;
    position: absolute;
    bottom: 10px;
    right: 10px;
    width: 25px;
    height: 25px;
    background-image: linear-gradient(
      135deg,
      ${n.borderLightest} 16.67%,
      ${n.material} 16.67%,
      ${n.material} 33.33%,
      ${n.borderDark} 33.33%,
      ${n.borderDark} 50%,
      ${n.borderLightest} 50%,
      ${n.borderLightest} 66.67%,
      ${n.material} 66.67%,
      ${n.material} 83.33%,
      ${n.borderDark} 83.33%,
      ${n.borderDark} 100%
    );
    background-size: 8.49px 8.49px;
    clip-path: polygon(100% 0px, 0px 100%, 100% 100%);
    cursor: nwse-resize;
  `}
`,Ws=m.forwardRef(({children:n,resizable:l=!1,resizeRef:s,shadow:u=!0,...f},g)=>B.createElement(Yg,{ref:g,shadow:u,...f},n,l&&B.createElement(qg,{"data-testid":"resizeHandle",ref:s})));Ws.displayName="Window";const Xg=ee(wr)`
  width: 234px;
  margin: 1rem 0;
  background: ${({theme:n})=>n.canvas};
`,Zg=ee.div`
  display: flex;
  background: ${({theme:n})=>n.materialDark};
  color: #dfe0e3;
`,Jg=ee.div`
  display: flex;
  flex-wrap: wrap;
`,xr=ee.div`
  text-align: center;
  height: 1.5em;
  line-height: 1.5em;
  width: 14.28%;
`,ev=ee.span`
  cursor: pointer;

  background: ${({active:n,theme:l})=>n?l.hoverBackground:"transparent"};
  color: ${({active:n,theme:l})=>n?l.canvasTextInvert:l.canvasText};

  &:hover {
    border: 2px dashed
      ${({theme:n,active:l})=>l?"none":n.materialDark};
  }
`,tv=[{value:0,label:"January"},{value:1,label:"February"},{value:2,label:"March"},{value:3,label:"April"},{value:4,label:"May"},{value:5,label:"June"},{value:6,label:"July"},{value:7,label:"August"},{value:8,label:"September"},{value:9,label:"October"},{value:10,label:"November"},{value:11,label:"December"}];function rv(n,l){return new Date(n,l+1,0).getDate()}function nv(n,l,s){return new Date(n,l,s).getDay()}function ov(n){const l=new Date(Date.parse(n)),s=l.getUTCDate(),u=l.getUTCMonth(),f=l.getUTCFullYear();return{day:s,month:u,year:f}}const iv=m.forwardRef(({className:n,date:l=new Date().toISOString(),onAccept:s,onCancel:u,shadow:f=!0},g)=>{const[h,v]=m.useState(()=>ov(l)),{year:x,month:E,day:T}=h,A=m.useCallback(({value:$})=>{v(Y=>({...Y,month:$}))},[]),P=m.useCallback($=>{v(Y=>({...Y,year:$}))},[]),U=m.useCallback($=>{v(Y=>({...Y,day:$}))},[]),j=m.useCallback(()=>{const $=[h.year,h.month+1,h.day].map(Y=>String(Y).padStart(2,"0")).join("-");s==null||s($)},[h.day,h.month,h.year,s]),y=m.useMemo(()=>{const $=Array.from({length:42}),Y=nv(x,E,1);let H=T;const V=rv(x,E);return H=H<V?H:V,$.forEach((J,L)=>{if(L>=Y&&L<V+Y){const M=L-Y+1;$[L]=B.createElement(xr,{key:L,onClick:()=>{U(M)}},B.createElement(ev,{active:M===H},M))}else $[L]=B.createElement(xr,{key:L})}),$},[T,U,E,x]);return B.createElement(Ws,{className:n,ref:g,shadow:f,style:{margin:20}},B.createElement(Fs,null,B.createElement("span",{role:"img","aria-label":"📆"},"📆"),"Date"),B.createElement(Bs,null,B.createElement($o,{noPadding:!0,style:{justifyContent:"space-between"}},B.createElement(Nn,{options:tv,value:E,onChange:A,width:128,menuMaxHeight:200}),B.createElement(Bf,{value:x,onChange:P,width:100})),B.createElement(Xg,null,B.createElement(Zg,null,B.createElement(xr,null,"S"),B.createElement(xr,null,"M"),B.createElement(xr,null,"T"),B.createElement(xr,null,"W"),B.createElement(xr,null,"T"),B.createElement(xr,null,"F"),B.createElement(xr,null,"S")),B.createElement(Jg,null,y)),B.createElement($o,{noPadding:!0,style:{justifyContent:"space-between"}},B.createElement(he,{fullWidth:!0,onClick:u,disabled:!u},"Cancel"),B.createElement(he,{fullWidth:!0,onClick:s?j:void 0,disabled:!s},"OK"))))});iv.displayName="DatePicker";const lv=n=>{switch(n){case"status":case"well":return me`
        ${nt({style:"status"})}
      `;case"window":case"outside":return me`
        ${nt({style:"window"})}
      `;case"field":return me`
        ${nt({style:"field"})}
      `;default:return me`
        ${nt()}
      `}},av=ee.div`
  position: relative;
  font-size: 1rem;
  ${({variant:n})=>lv(n)}
  ${({variant:n})=>er(n==="field"?{background:"canvas",color:"canvasText"}:void 0)}
`,an=m.forwardRef(({children:n,shadow:l=!1,variant:s="window",...u},f)=>B.createElement(av,{ref:f,shadow:l,variant:s,...u},n));an.displayName="Frame";const sv=ee.fieldset`
  position: relative;
  border: 2px solid
    ${({theme:n,variant:l})=>l==="flat"?n.flatDark:n.borderLightest};
  padding: 16px;
  margin-top: 8px;
  font-size: 1rem;
  color: ${({theme:n})=>n.materialText};
  ${({variant:n})=>n!=="flat"&&me`
      box-shadow: -1px -1px 0 1px ${({theme:l})=>l.borderDark},
        inset -1px -1px 0 1px ${({theme:l})=>l.borderDark};
    `}
  ${n=>n.$disabled&&Jt()}
`,cv=ee.legend`
  display: flex;
  position: absolute;
  top: 0;
  left: 8px;
  transform: translateY(calc(-50% - 2px));
  padding: 0 8px;

  font-size: 1rem;
  background: ${({theme:n,variant:l})=>l==="flat"?n.canvas:n.material};
`,Ct=m.forwardRef(({label:n,disabled:l=!1,variant:s="default",children:u,...f},g)=>B.createElement(sv,{"aria-disabled":l,$disabled:l,variant:s,ref:g,...f},n&&B.createElement(cv,{variant:s},n),u));Ct.displayName="GroupBox";const uv=ee.div`
  ${({theme:n,size:l="100%"})=>`
  display: inline-block;
  box-sizing: border-box;
  height: ${Fr(l)};
  width: 5px;
  border-top: 2px solid ${n.borderLightest};
  border-left: 2px solid ${n.borderLightest};
  border-bottom: 2px solid ${n.borderDark};
  border-right: 2px solid ${n.borderDark};
  background: ${n.material};
`}
`;uv.displayName="Handle";const dv="url('data:image/gif;base64,R0lGODlhPAA8APQAADc3N6+vr4+Pj05OTvn5+V1dXZ+fn29vby8vLw8PD/X19d/f37S0tJSUlLq6und3d39/f9XV1c/Pz+bm5qamphkZGWZmZsbGxr+/v+rq6tra2u/v7yIiIv///wAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh+QQFBAAfACH+I1Jlc2l6ZWQgb24gaHR0cHM6Ly9lemdpZi5jb20vcmVzaXplACwAAAAAPAA8AAAF/+AnjmRpnmiqrmzrvnAsz3Rt37jr7Xzv/8BebhQsGn1D0XFZTH6YUGQySvU4fYKAdsvtdi1Cp3In6ZjP6HTawBMTyWbFYk6v18/snXvsKXciUApmeVZ7PH6ATIIdhHtPcB0TDQ1gQBCTBINthpBnAUEaa5tuh2mfQKFojZx9aRMSEhA7FLAbonqsfmoUOxFqmriknWm8Hr6/q8IeCAAAx2cTERG2aBTNHMGOj8a/v8WF2m/c3cSj4SQ8C92n4Ocm6evm7ui9CosdBPbs8yo8E2YO5PE74Q+gwIElCnYImA3hux3/Fh50yCciw3YUt2GQtiiDtGQO4f3al1GkGpIDeXlg0KDhXpoMLBtMVPaMnJlv/HjUtIkzHA8HEya4tLkhqICGV4bZVAMyaaul3ZpOUQoVz8wbpaoyvWojq1ZVXGt4/QoM49SnZMs6GktW6hC2X93mgKtVbtceWbzo9VIJKdYqUJwCPiJ4cJOzhg+/TWwko+PHkCNLdhgCACH5BAUEAB8ALAAAAAABAAEAAAUD4BcCACH5BAUEAB8ALBYADAAQAA0AAAVFYCeOZPmVaKqimeO+MPxFXv3d+F17Cm3nuJ1ic7lAdroapUjABZCfnQb4ef6k1OHGULtsNk3qjVKLiIFkj/mMIygU4VwIACH5BAUEAB8ALAAAAAABAAEAAAUD4BcCACH5BAUEAB8ALBkAIwAKAAcAAAUp4CdehrGI6Ed5XpSKa4teguBoGlVPAXuJBpam5/l9gh7NZrFQiDJMRQgAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsFgAPABAAIQAABVBgJ45kaZ5oakZB67bZ+M10bd94ru987//AoHBILNYYAsGlR/F4IkwnlLeZTBQ9UlaWwzweERHjuzAKFZkMYYZWm4mOw0ETfdanO8Vms7aFAAAh+QQFBAAfACwAAAAAAQABAAAFA+AXAgAh+QQFBAAfACwZABIACgAeAAAFUGAnjmRpnij5rerqtu4Hx3Rt33iu758iZrUZa1TDCASLGsXjiSiZzmFnM5n4TNJSdmREElfL5lO8cgwGACbgrAkwPat3+x1naggKRS+f/4QAACH5BAUEAB8ALAAAAAABAAEAAAUD4BcCACH5BAUEAB8ALBYAIwAQAA0AAAVE4CeOXdmNaGqeabu27SUIC5xSnifZKK7zl8djkCsIaylGziNaakaEzcbH/Cwl0k9kuWxyPYptzrZULA7otFpNIK1eoxAAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkEBQQAHwAsAAAAAAEAAQAABQPgFwIAIfkECQQAHwAsDgAEACAANAAABTHgJ45kaZ5oqq5s675wLM90bd94ru987//AoHBILBqPyKRyyWw6n9CodEqtWq/Y7CoEACH5BAUEAB8ALAAAAAA8ADwAAAX/4CeOZGmeaKqubFt6biy3Xj3fuFjveU/vPJ/wBAQOj6RiEClUGpk9IMAJxQEdmQK1Grt2OhutkvurOb7f8JaM8qLT4iKbuDu/0erxfOS+4+NPex9mfn55coIfCAuFhoBLbDUAjI1vh4FkOxSVd5eQXB4GnI5rXAAbo6R6VTUFqKmWjzasNaKwsaVIHhAEt3cLTjBQA6++XwoHuUM1vMYdyMorwoN8wkC2t9A8s102204Wxana3DNAAQO1FjUCEDXhvuTT5nUdEwOiGxa8BBDwXxKaLTiAKoMFRvJy9CmmoFcHAgrQSEiwKwICDwU0pAMQIdmnboR8TfwWrJyMPrAiz1DkNs2aSRbe6hnr99LEvDJ9IB5DQ8Dhm36glNh5COGBAmQNHrbz+WXBFChOTqFx5+GBxwYCmL1ZcPHmMiWuvkTgECzBBUvrvH4tErbDWCcYDB2IBPbV2yJJ72SZ46TtXSB5v2RIp1ZXXbFkgWxCc68mk752E3tY/OZeIsiIaxi9o+BBokGH3SZ+4FPbZ8yiPQxNeDl0hNUeHWcKjYb1Zx20bd/GzRaV7t28gRSYELvw7pIfgVcLplwF8+bOo0Ffjmm6zerWrxvPzoe79w8hAAAh+QQJBAAfACwBAAEAOgA6AAAFRuAnjmRpnmiqrmzrvnAsz3Rt33iu73zv/8CgcEgsGo/IpHLJbDqf0Kh0Sq1ar9isdsvter/gsHhMLpvP6LR6zW673/D4MgQAIfkEBQQAHwAsAAAAADwAPAAABf/gJ45kaZ5oqq5s675wLM90bd94ru987//AoHBILBqPyJxnyTQym6nn0ilVSa9XGHY7jXKx2m/WK36Gy1CUVCBpu9+OtNqDeNslgip5Gej4/4ATcidLAICHHQF6c0x9iH+CXV6Gj36KZnsejgsREQSACp0Yg0ydEZWWi4RPjgdLG48apEuogJeDJVKtr7GzHrV/t5KrjX6uHhQMF4cKCwujTxHOwKmYjHzGTw+VEVIK1MGqJrrZTNuP3U/f4IniuazlSwMUFMugE/j47NW4JOQdx9bsoybMgxV4ALEIGAis4MFiCZkUaLPgUAYHGDF+Yucw0y5z3Lzt63hNUzwP5xCRpWOyDhxJYtgiStBQEVCGAAEM6MLp0p0/hMdgIZI17AOTntZgmowo9BBRgz9/EfQ54h8BBS39bKDXwBc9CrVejkNYKRLUSWGpivhXtt9PSpXEvmNiwYDdu3jzFB3LAa9fAxbUGkXjtmSZh4TPJM4kRgbhvVEL9xhTEongJJgza97MubPnz6BDix5NurTp0yJCAAAh+QQJBAAfACwEAA4ANAAgAAAFMeAnjmRpnmiqrmzrvnAsz3Rt33iu73zv/8CgcEgsGo/IpHLJbDqf0Kh0Sq1ar9jsKgQAIfkEBQQAHwAsAAAAADwAPAAABf/gJ45kaZ5oqq5s6bVwLHu0bN8uXeM8rP+9YOoHFBpHRN1xmSwue02A82lrFjaOKbVl3XQ6WeWWm7x+v+HdeFj2ntHaNbL9jUAI5/RLTurWOR53eXFbfh0RgB4PCm9hfCKGiDSLb18Bjx+RiR4HjG8TA3trmkSdZxuhalSkRA2VBqpPrD+ulR0Go3SHmz8CeG8bFqJMupJNHr5nCsKxQccTg4oUNA0YCYG/HQQQYsSlnmCUFLUXgm8EAsPeP6Zf2baV2+rEmTrt8PDyzS7O9uD4b5YV2VGjGw52/wB+CaYjlQcpNBAQioHwy4QMCxe4i3BKGIQN3K7AArBATz8anUDADcgQDMGCbQkknDKAh4ABNxQ0gpnoQ8eDVAUO0ADAzUNMhbZMQiG4R4mOo0gb8eTCQgeEqJVM7juCDWvWJnI4ev2aZIwHl2PfZIBIZBXKtAsLgC1kJu0GuWXNaoB7d67ZlWP75jVLw4JXwW35PNSJFPFUrmIb402smFNCW44N5kJ5+dTkx+vuAfus+VHF0X4xzeHsObXq1ZY7ZN76mt0C0rRf1zuWW/du175PHAu+YjhxFcCPm6CsHHnv5kig6w4BACH5BAkEAB8ALAEAAQA6ADoAAAVG4CeOZGmeaKqubOu+cCzPdG3feK7vfO//wKBwSCwaj8ikcslsOp/QqHRKrVqv2Kx2y+16v+CweEwum8/otHrNbrvf8PgyBAAh+QQFBAAfACwAAAAAPAA8AAAF/+AnjmRpnmiqrmzrvnAsz3Rt37jr7Xzv/8BebhQsGn1D0XFZTH6YUGQySvU4fYKAdsvtdi1Cp3In6ZjP6HTawBMTyWbFYk6v18/snXvsKXciUApmeVZ7PH6ATIIdhHtPcB0TDQ1gQBCTBINthpBnAUEaa5tuh2mfQKFojZx9aRMSEhA7FLAbonqsfmoUOxFqmriknWm8Hr6/q8IeCAAAx2cTERG2aBTNHMGOj8a/v8WF2m/c3cSj4SQ8C92n4Ocm6evm7ui9CosdBPbs8yo8E2YO5PE74Q+gwIElCnYImA3hux3/Fh50yCciw3YUt2GQtiiDtGQO4f3al1GkGpIDeXlg0KDhXpoMLBtMVPaMnJlv/HjUtIkzHA8HEya4tLkhqICGV4bZVAMyaaul3ZpOUQoVz8wbpaoyvWojq1ZVXGt4/QoM49SnZMs6GktW6hC2X93mgKtVbtceWbzo9VIJKdYqUJwCPiJ4cJOzhg+/TWwko+PHkCNLdhgCACH5BAUEAB8ALAAAAAABAAEAAAUD4BcCADs=')",fv=ee.div`
  display: inline-block;
  height: ${({size:n})=>Fr(n)};
  width: ${({size:n})=>Fr(n)};
`,pv=ee.span`
  display: block;
  background: ${dv};
  background-size: cover;
  width: 100%;
  height: 100%;
`,hv=m.forwardRef(({size:n=30,...l},s)=>B.createElement(fv,{size:n,ref:s,...l},B.createElement(pv,null)));hv.displayName="Hourglass";const mv=ee.div`
  position: relative;
  display: inline-block;
  padding-bottom: 26px;
`,gv=ee.div`
  position: relative;
`,vv=ee.div`
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  width: 195px;
  height: 155px;
  padding: 12px;
  background: ${({theme:n})=>n.material};
  border-top: 4px solid ${({theme:n})=>n.borderLightest};
  border-left: 4px solid ${({theme:n})=>n.borderLightest};
  border-bottom: 4px solid ${({theme:n})=>n.borderDark};
  border-right: 4px solid ${({theme:n})=>n.borderDark};

  outline: 1px dotted ${({theme:n})=>n.material};
  outline-offset: -3px;
  &:before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    outline: 1px dotted ${({theme:n})=>n.material};
  }
  box-shadow: 1px 1px 0 1px ${({theme:n})=>n.borderDarkest};

  &:after {
    content: '';
    display: inline-block;
    position: absolute;
    bottom: 4px;
    right: 12px;
    width: 10px;
    border-top: 2px solid #4d9046;
    border-bottom: 2px solid #07ff00;
  }
`,yv=ee(cr).attrs(()=>({"data-testid":"background"}))`
  width: 100%;
  height: 100%;
`,xv=ee.div`
  box-sizing: border-box;
  position: absolute;
  top: calc(100% + 2px);
  left: 50%;
  transform: translateX(-50%);
  height: 10px;
  width: 50%;
  background: ${({theme:n})=>n.material};
  border-left: 2px solid ${({theme:n})=>n.borderLightest};
  border-bottom: 2px solid ${({theme:n})=>n.borderDarkest};
  border-right: 2px solid ${({theme:n})=>n.borderDarkest};
  box-shadow: inset 0px 0px 0px 2px ${({theme:n})=>n.borderDark};

  &:before {
    content: '';
    position: absolute;
    top: calc(100% + 2px);
    left: 50%;
    transform: translateX(-50%);
    width: 80%;
    height: 8px;
    background: ${({theme:n})=>n.material};
    border-left: 2px solid ${({theme:n})=>n.borderLightest};
    border-right: 2px solid ${({theme:n})=>n.borderDarkest};
    box-shadow: inset 0px 0px 0px 2px ${({theme:n})=>n.borderDark};
  }
  &:after {
    content: '';
    position: absolute;
    top: calc(100% + 8px);
    left: 50%;
    transform: translateX(-50%);
    width: 150%;
    height: 4px;
    background: ${({theme:n})=>n.material};
    border: 2px solid ${({theme:n})=>n.borderDark};
    border-bottom: none;
    box-shadow: inset 1px 1px 0px 1px ${({theme:n})=>n.borderLightest},
      1px 1px 0 1px ${({theme:n})=>n.borderDarkest};
  }
`,wv=m.forwardRef(({backgroundStyles:n,children:l,...s},u)=>B.createElement(mv,{ref:u,...s},B.createElement(gv,null,B.createElement(vv,null,B.createElement(yv,{style:n},l)),B.createElement(xv,null))));wv.displayName="Monitor";const bv=ee.div`
  display: inline-block;
  height: ${gt.md};
  width: 100%;
`,kv=ee(cr)`
  width: 100%;
  height: 100%;
  position: relative;
  text-align: center;
  padding: 0;
  overflow: hidden;
  &:before {
    z-index: 1;
  }
`,Kf=me`
  width: calc(100% - 4px);
  height: calc(100% - 4px);

  display: flex;
  align-items: center;
  justify-content: space-around;
`,Sv=ee.div`
  position: relative;
  top: 4px;
  ${Kf}
  background: ${({theme:n})=>n.canvas};
  color: #000;
  margin-left: 2px;
  margin-top: -2px;
  color: ${({theme:n})=>n.materialText};
`,Cv=ee.div`
  position: absolute;
  top: 2px;
  left: 2px;
  ${Kf}
  color: ${({theme:n})=>n.materialTextInvert};
  background: ${({theme:n})=>n.progress};
  clip-path: polygon(
    0 0,
    ${({value:n=0})=>n}% 0,
    ${({value:n=0})=>n}% 100%,
    0 100%
  );
  transition: 0.4s linear clip-path;
`,Av=ee.div`
  width: calc(100% - 6px);
  height: calc(100% - 8px);
  position: absolute;
  left: 3px;
  top: 4px;
  box-sizing: border-box;
  display: inline-flex;
`,Yf=17,Ev=ee.span`
  display: inline-block;
  width: ${Yf}px;
  box-sizing: border-box;
  height: 100%;
  background: ${({theme:n})=>n.progress};
  border-color: ${({theme:n})=>n.material};
  border-width: 0px 1px;
  border-style: solid;
`,_v=m.forwardRef(({hideValue:n=!1,shadow:l=!0,value:s,variant:u="default",...f},g)=>{const h=n?null:`${s}%`,v=m.useRef(null),[x,E]=m.useState([]),T=m.useCallback(()=>{if(!v.current||s===void 0)return;const A=v.current.getBoundingClientRect().width,P=Math.round(s/100*A/Yf);E(Array.from({length:P}))},[s]);return m.useEffect(()=>(T(),window.addEventListener("resize",T),()=>window.removeEventListener("resize",T)),[T]),B.createElement(bv,{"aria-valuenow":s!==void 0?Math.round(s):void 0,ref:g,role:"progressbar",variant:u,...f},B.createElement(kv,{variant:u,shadow:l},u==="default"?B.createElement(B.Fragment,null,B.createElement(Sv,{"data-testid":"defaultProgress1"},h),B.createElement(Cv,{"data-testid":"defaultProgress2",value:s},h)):B.createElement(Av,{ref:v,"data-testid":"tileProgress"},x.map((A,P)=>B.createElement(Ev,{key:P})))))});_v.displayName="ProgressBar";const qf=me`
  width: ${sr}px;
  height: ${sr}px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin-right: 0.5rem;
`,Tv=ee(cr)`
  ${qf}
  background: ${({$disabled:n,theme:l})=>n?l.material:l.canvas};

  &:before {
    content: '';
    position: absolute;
    left: 0px;
    top: 0px;
    width: calc(100% - 4px);
    height: calc(100% - 4px);
    border-radius: 50%;
    box-shadow: none;
  }
`,jv=ee.div`
  ${un()}
  ${qf}
  outline: none;
  background: ${({$disabled:n,theme:l})=>n?l.flatLight:l.canvas};
  &:before {
    content: '';
    display: inline-block;
    position: absolute;
    top: 0;
    left: 0;
    width: calc(100% - 4px);
    height: calc(100% - 4px);
    border: 2px solid ${({theme:n})=>n.flatDark};
    border-radius: 50%;
  }
`,Pv=ee.span.attrs(()=>({"data-testid":"checkmarkIcon"}))`
  position: absolute;
  content: '';
  display: inline-block;
  top: 50%;
  left: 50%;
  width: 6px;
  height: 6px;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  background: ${n=>n.$disabled?n.theme.checkmarkDisabled:n.theme.checkmark};
`,Rv={flat:jv,default:Tv},$v=m.forwardRef(({checked:n,className:l="",disabled:s=!1,label:u="",onChange:f,style:g={},variant:h="default",...v},x)=>{const E=Rv[h];return B.createElement(Ms,{$disabled:s,className:l,style:g},B.createElement(E,{$disabled:s,role:"presentation"},n&&B.createElement(Pv,{$disabled:s,variant:h})),B.createElement(rl,{disabled:s,onChange:s?void 0:f,readOnly:s,type:"radio",checked:n,ref:x,...v}),u&&B.createElement(Os,null,u))});$v.displayName="Radio";const Dv=typeof window<"u"?m.useLayoutEffect:m.useEffect;function nn(n){const l=m.useRef(n);return Dv(()=>{l.current=n}),m.useCallback((...s)=>(0,l.current)(...s),[])}function ef(n,l){typeof n=="function"?n(l):n&&(n.current=l)}function tf(n,l){return m.useMemo(()=>n==null&&l==null?null:s=>{ef(n,s),ef(l,s)},[n,l])}var zv=kf();let ol=!0,Ts=!1,rf;const Iv={text:!0,search:!0,url:!0,tel:!0,email:!0,password:!0,number:!0,date:!0,month:!0,week:!0,time:!0,datetime:!0,"datetime-local":!0};function Lv(n){if("type"in n){const{type:l,tagName:s}=n;if(s==="INPUT"&&Iv[l]&&!n.readOnly||s==="TEXTAREA"&&!n.readOnly)return!0}return!!("isContentEditable"in n&&n.isContentEditable)}function Mv(n){n.metaKey||n.altKey||n.ctrlKey||(ol=!0)}function is(){ol=!1}function Ov(){this.visibilityState==="hidden"&&Ts&&(ol=!0)}function Nv(n){n.addEventListener("keydown",Mv,!0),n.addEventListener("mousedown",is,!0),n.addEventListener("pointerdown",is,!0),n.addEventListener("touchstart",is,!0),n.addEventListener("visibilitychange",Ov,!0)}function Bv(n){const{target:l}=n;try{return l.matches(":focus-visible")}catch{}return ol||Lv(l)}function Fv(){Ts=!0,window.clearTimeout(rf),rf=window.setTimeout(()=>{Ts=!1},100)}function Wv(){const n=m.useCallback(l=>{const s=zv.findDOMNode(l);s!=null&&Nv(s.ownerDocument)},[]);return{isFocusVisible:Bv,onBlurVisible:Fv,ref:n}}function Hv(n,l,s){return(s-l)*n+l}function Ki(n,l){if(l!==void 0&&"changedTouches"in n){for(let s=0;s<n.changedTouches.length;s+=1){const u=n.changedTouches[s];if(u.identifier===l)return{x:u.clientX,y:u.clientY}}return!1}return"clientX"in n?{x:n.clientX,y:n.clientY}:!1}function Yi(n){return n&&n.ownerDocument||document}function Uv(n,l){var s;const{index:u}=(s=n.reduce((f,g,h)=>{const v=Math.abs(l-g);return f===null||v<f.distance||v===f.distance?{distance:v,index:h}:f},null))!==null&&s!==void 0?s:{};return u??-1}const Vv=ee.div`
  display: inline-block;
  position: relative;
  touch-action: none;
  &:before {
    content: '';
    display: inline-block;
    position: absolute;
    top: -2px;
    left: -15px;
    width: calc(100% + 30px);
    height: ${({hasMarks:n})=>n?"41px":"39px"};
    ${({isFocused:n,theme:l})=>n&&`
        outline: 2px dotted ${l.materialText};
        `}
  }

  ${({orientation:n,size:l})=>n==="vertical"?me`
          height: ${l};
          margin-right: 1.5rem;
          &:before {
            left: -6px;
            top: -15px;
            height: calc(100% + 30px);
            width: ${({hasMarks:s})=>s?"41px":"39px"};
          }
        `:me`
          width: ${l};
          margin-bottom: 1.5rem;
          &:before {
            top: -2px;
            left: -15px;
            width: calc(100% + 30px);
            height: ${({hasMarks:s})=>s?"41px":"39px"};
          }
        `}

  pointer-events: ${({$disabled:n})=>n?"none":"auto"};
`,Xf=()=>me`
  position: absolute;
  ${({orientation:n})=>n==="vertical"?me`
          bottom: 0;
          left: 50%;
          transform: translateX(-50%);
          height: 100%;
          width: 8px;
        `:me`
          left: 0;
          top: 50%;
          transform: translateY(-50%);
          height: 8px;
          width: 100%;
        `}
`,Qv=ee(cr)`
  ${Xf()}
`,Gv=ee(cr)`
  ${Xf()}

  border-left-color: ${({theme:n})=>n.flatLight};
  border-top-color: ${({theme:n})=>n.flatLight};
  border-right-color: ${({theme:n})=>n.canvas};
  border-bottom-color: ${({theme:n})=>n.canvas};
  &:before {
    border-left-color: ${({theme:n})=>n.flatDark};
    border-top-color: ${({theme:n})=>n.flatDark};
    border-right-color: ${({theme:n})=>n.flatLight};
    border-bottom-color: ${({theme:n})=>n.flatLight};
  }
`,Kv=ee.span`
  position: relative;
  ${({orientation:n})=>n==="vertical"?me`
          width: 32px;
          height: 18px;
          right: 2px;
          transform: translateY(-50%);
        `:me`
          height: 32px;
          width: 18px;
          top: 2px;
          transform: translateX(-50%);
        `}
  ${({variant:n})=>n==="flat"?me`
          ${un()}
          outline: 2px solid ${({theme:l})=>l.flatDark};
          background: ${({theme:l})=>l.flatLight};
        `:me`
          ${er()}
          ${nt()}
          &:focus {
            outline: none;
          }
        `}
    ${({$disabled:n,theme:l})=>n&&Do({mainColor:l.material,secondaryColor:l.borderLightest})}
`,On=6,Yv=ee.span`
  display: inline-block;
  position: absolute;

  ${({orientation:n})=>n==="vertical"?me`
          right: ${-On-2}px;
          bottom: 0px;
          transform: translateY(1px);
          width: ${On}px;
          border-bottom: 2px solid ${({theme:l})=>l.materialText};
        `:me`
          bottom: ${-On}px;
          height: ${On}px;
          transform: translateX(-1px);
          border-left: 1px solid ${({theme:l})=>l.materialText};
          border-right: 1px solid ${({theme:l})=>l.materialText};
        `}

  color:  ${({theme:n})=>n.materialText};
  ${({$disabled:n,theme:l})=>n&&me`
      ${Jt()}
      box-shadow: 1px 1px 0px ${l.materialTextDisabledShadow};
      border-color: ${l.materialTextDisabled};
    `}
`,qv=ee.div`
  position: absolute;
  bottom: 0;
  left: 0;
  line-height: 1;
  font-size: 0.875rem;

  ${({orientation:n})=>n==="vertical"?me`
          transform: translate(${On+2}px, ${On+1}px);
        `:me`
          transform: translate(-0.5ch, calc(100% + 2px));
        `}
`,Xv=m.forwardRef(({defaultValue:n,disabled:l=!1,marks:s=!1,max:u=100,min:f=0,name:g,onChange:h,onChangeCommitted:v,onMouseDown:x,orientation:E="horizontal",size:T="100%",step:A=1,value:P,variant:U="default",...j},y)=>{const $=U==="flat"?Gv:Qv,Y=E==="vertical",[H=f,V]=Wr({defaultValue:n,onChange:h??v,value:P}),{isFocusVisible:J,onBlurVisible:L,ref:M}=Wv(),[ne,de]=m.useState(!1),Ae=m.useRef(),$e=m.useRef(null),De=tf(M,Ae),W=tf(y,De),ce=nn(I=>{J(I)&&de(!0)}),xe=nn(()=>{ne!==!1&&(de(!1),L())}),K=m.useRef(),ue=m.useMemo(()=>s===!0&&Number.isFinite(A)?[...Array(Math.round((u-f)/A)+1)].map((I,D)=>({label:void 0,value:f+A*D})):Array.isArray(s)?s:[],[s,u,f,A]),ge=nn(I=>{const D=(u-f)/10,q=ue.map(Q=>Q.value),b=q.indexOf(H);let oe=0;switch(I.key){case"Home":oe=f;break;case"End":oe=u;break;case"PageUp":A&&(oe=H+D);break;case"PageDown":A&&(oe=H-D);break;case"ArrowRight":case"ArrowUp":A?oe=H+A:oe=q[b+1]||q[q.length-1];break;case"ArrowLeft":case"ArrowDown":A?oe=H-A:oe=q[b-1]||q[0];break;default:return}I.preventDefault(),A&&(oe=Xd(oe,A,f)),oe=ln(oe,f,u),V(oe),de(!0),h==null||h(oe),v==null||v(oe)}),F=m.useCallback(I=>{if(!Ae.current)return 0;const D=Ae.current.getBoundingClientRect();let q;Y?q=(D.bottom-I.y)/D.height:q=(I.x-D.left)/D.width;let b;if(b=Hv(q,f,u),A)b=Xd(b,A,f);else{const oe=ue.map(_=>_.value),Q=Uv(oe,b);b=oe[Q]}return b=ln(b,f,u),b},[ue,u,f,A,Y]),se=nn(I=>{var D;const q=Ki(I,K.current);if(!q)return;const b=F(q);(D=$e.current)===null||D===void 0||D.focus(),V(b),de(!0),h==null||h(b)}),G=nn(I=>{const D=Ki(I,K.current);if(!D)return;const q=F(D);v==null||v(q),K.current=void 0;const b=Yi(Ae.current);b.removeEventListener("mousemove",se),b.removeEventListener("mouseup",G),b.removeEventListener("touchmove",se),b.removeEventListener("touchend",G)}),S=nn(I=>{var D;x==null||x(I),I.preventDefault(),(D=$e.current)===null||D===void 0||D.focus(),de(!0);const q=Ki(I,K.current);if(q){const oe=F(q);V(oe),h==null||h(oe)}const b=Yi(Ae.current);b.addEventListener("mousemove",se),b.addEventListener("mouseup",G)}),C=nn(I=>{var D;I.preventDefault();const q=I.changedTouches[0];q!=null&&(K.current=q.identifier),(D=$e.current)===null||D===void 0||D.focus(),de(!0);const b=Ki(I,K.current);if(b){const Q=F(b);V(Q),h==null||h(Q)}const oe=Yi(Ae.current);oe.addEventListener("touchmove",se),oe.addEventListener("touchend",G)});return m.useEffect(()=>{const{current:I}=Ae;I==null||I.addEventListener("touchstart",C);const D=Yi(I);return()=>{I==null||I.removeEventListener("touchstart",C),D.removeEventListener("mousemove",se),D.removeEventListener("mouseup",G),D.removeEventListener("touchmove",se),D.removeEventListener("touchend",G)}},[G,se,C]),B.createElement(Vv,{$disabled:l,hasMarks:!!ue.length,isFocused:ne,onMouseDown:S,orientation:E,ref:W,size:Fr(T),...j},B.createElement("input",{disabled:l,name:g,type:"hidden",value:H??0}),ue&&ue.map(I=>B.createElement(Yv,{$disabled:l,"data-testid":"tick",key:I.value/(u-f)*100,orientation:E,style:{[Y?"bottom":"left"]:`${(I.value-f)/(u-f)*100}%`}},I.label&&B.createElement(qv,{"aria-hidden":!0,"data-testid":"mark",orientation:E},I.label))),B.createElement($,{orientation:E,variant:U}),B.createElement(Kv,{$disabled:l,"aria-disabled":l?!0:void 0,"aria-orientation":E,"aria-valuemax":u,"aria-valuemin":f,"aria-valuenow":H,onBlur:xe,onFocus:ce,onKeyDown:ge,orientation:E,ref:$e,role:"slider",style:{[Y?"bottom":"left"]:`${(Y?-100:0)+100*(H-f)/(u-f)}%`},tabIndex:l?void 0:0,variant:U}))});Xv.displayName="Slider";const Zv=ee.tbody`
  background: ${({theme:n})=>n.canvas};
  display: table-row-group;
  box-shadow: ${Is};
  overflow-y: auto;
`,zo=m.forwardRef(function({children:l,...s},u){return B.createElement(Zv,{ref:u,...s},l)});zo.displayName="TableBody";const Jv=ee.td`
  padding: 0 8px;
`,et=m.forwardRef(function({children:l,...s},u){return B.createElement(Jv,{ref:u,...s},l)});et.displayName="TableDataCell";const e0=ee.thead`
  display: table-header-group;
`,Io=m.forwardRef(function({children:l,...s},u){return B.createElement(e0,{ref:u,...s},l)});Io.displayName="TableHead";const t0=ee.th`
  position: relative;
  padding: 0 8px;
  display: table-cell;
  vertical-align: inherit;
  background: ${({theme:n})=>n.material};
  cursor: default;
  user-select: none;
  &:before {
    box-sizing: border-box;
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    ${nt()}

    border-left: none;
    border-top: none;
  }
  ${({$disabled:n})=>!n&&me`
      &:active {
        &:before {
          ${nt({invert:!0,style:"window"})}
          border-left: none;
          border-top: none;
          padding-top: 2px;
        }

        & > div {
          position: relative;
          top: 2px;
        }
      }
    `}

  color: ${({theme:n})=>n.materialText};
  ${({$disabled:n})=>n&&Jt()}
  &:hover {
    color: ${({theme:n})=>n.materialText};
    ${({$disabled:n})=>n&&Jt()}
  }
`,tt=m.forwardRef(function({disabled:l=!1,children:s,onClick:u,onTouchStart:f=dn,sort:g,...h},v){const x=g==="asc"?"ascending":g==="desc"?"descending":void 0;return B.createElement(t0,{$disabled:l,"aria-disabled":l,"aria-sort":x,onClick:l?void 0:u,onTouchStart:l?void 0:f,ref:v,...h},B.createElement("div",null,s))});tt.displayName="TableHeadCell";const r0=ee.tr`
  color: inherit;
  display: table-row;
  height: calc(${gt.md} - 2px);
  line-height: calc(${gt.md} - 2px);
  vertical-align: middle;
  outline: none;

  color: ${({theme:n})=>n.canvasText};
  &:hover {
    background: ${({theme:n})=>n.hoverBackground};
    color: ${({theme:n})=>n.canvasTextInvert};
  }
`,Ht=m.forwardRef(function({children:l,...s},u){return B.createElement(r0,{ref:u,...s},l)});Ht.displayName="TableRow";const n0=ee.table`
  display: table;
  width: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  font-size: 1rem;
`,o0=ee(cr)`
  &:before {
    box-shadow: none;
  }
`,Lo=m.forwardRef(({children:n,...l},s)=>B.createElement(o0,null,B.createElement(n0,{ref:s,...l},n)));Lo.displayName="Table";const i0=ee.button`
  ${er()}
  ${nt()}
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  height: ${gt.md};
  line-height: ${gt.md};
  padding: 0 8px;
  border-bottom: none;
  border-top-left-radius: 5px;
  border-top-right-radius: 5px;
  margin: 0 0 -2px 0;
  cursor: default;
  color: ${({theme:n})=>n.materialText};
  user-select: none;
  font-family: inherit;
  &:focus:after,
  &:active:after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    ${Wn}
    outline-offset: -6px;
  }
  ${n=>n.selected&&`
    z-index: 1;
    height: calc(${gt.md} + 4px);
    top: -4px;
    margin-bottom: -6px;
    padding: 0 16px;
    margin-left: -8px;
    &:not(:last-child) {
      margin-right: -8px;
    }
  `}
  &:before {
    content: '';
    position: absolute;
    width: calc(100% - 4px);
    height: 6px;
    background: ${({theme:n})=>n.material};
    bottom: -4px;
    left: 2px;
  }
`,ar=m.forwardRef(({value:n,onClick:l,selected:s=!1,children:u,...f},g)=>B.createElement(i0,{"aria-selected":s,selected:s,onClick:h=>l==null?void 0:l(n,h),ref:g,role:"tab",...f},u));ar.displayName="Tab";const l0=ee.div`
  ${er()}
  ${nt()}
  position: relative;
  display: block;
  height: 100%;
  padding: 16px;
  font-size: 1rem;
`,il=m.forwardRef(({children:n,...l},s)=>B.createElement(l0,{ref:s,...l},n));il.displayName="TabBody";const a0=ee.div`
  position: relative;
  ${({isMultiRow:n,theme:l})=>n&&`
  button {
    flex-grow: 1;
  }
  button:last-child:before {
    border-right: 2px solid ${l.borderDark};
  }
  `}
`,s0=ee.div.attrs(()=>({"data-testid":"tab-row"}))`
  position: relative;
  display: flex;
  flex-wrap: no-wrap;
  text-align: left;
  left: 8px;
  width: calc(100% - 8px);

  &:not(:first-child):before {
    content: '';
    position: absolute;
    right: 0;
    left: 0;
    height: 100%;
    border-right: 2px solid ${({theme:n})=>n.borderDarkest};
    border-left: 2px solid ${({theme:n})=>n.borderLightest};
  }
`;function c0(n,l){const s=[];for(let u=l;u>0;u-=1)s.push(n.splice(0,Math.ceil(n.length/u)));return s}const ll=m.forwardRef(({value:n,onChange:l=dn,children:s,rows:u=1,...f},g)=>{const h=m.useMemo(()=>{var v;const x=(v=B.Children.map(s,A=>{if(!B.isValidElement(A))return null;const P={selected:A.props.value===n,onClick:l};return B.cloneElement(A,P)}))!==null&&v!==void 0?v:[],E=c0(x,u).map((A,P)=>({key:P,tabs:A})),T=E.findIndex(A=>A.tabs.some(P=>P.props.selected));return E.push(E.splice(T,1)[0]),E},[s,l,u,n]);return B.createElement(a0,{...f,isMultiRow:u>1,role:"tablist",ref:g},h.map(v=>B.createElement(s0,{key:v.key},v.tabs)))});ll.displayName="Tabs";const u0=["blur","focus"],d0=["click","contextmenu","doubleclick","drag","dragend","dragenter","dragexit","dragleave","dragover","dragstart","drop","mousedown","mouseenter","mouseleave","mousemove","mouseout","mouseover","mouseup"];function nf(n){return"nativeEvent"in n&&u0.includes(n.type)}function of(n){return"nativeEvent"in n&&d0.includes(n.type)}const f0={top:`top: -4px;
        left: 50%;
        transform: translate(-50%, -100%);`,bottom:`bottom: -4px;
           left: 50%;
           transform: translate(-50%, 100%);`,left:`left: -4px;
         top: 50%;
         transform: translate(-100%, -50%);`,right:`right: -4px;
          top: 50%;
          transform: translate(100%, -50%);`},p0=ee.span`
  position: absolute;

  z-index: 1;
  display: ${n=>n.show?"block":"none"};
  padding: 4px;
  border: 2px solid ${({theme:n})=>n.borderDarkest};
  background: ${({theme:n})=>n.tooltip};
  box-shadow: ${zs};
  text-align: center;
  font-size: 1rem;
  ${n=>f0[n.position]}
`,h0=ee.div`
  position: relative;
  display: inline-block;
  white-space: nowrap;
`,m0=m.forwardRef(({className:n,children:l,disableFocusListener:s=!1,disableMouseListener:u=!1,enterDelay:f=1e3,leaveDelay:g=0,onBlur:h,onClose:v,onFocus:x,onMouseEnter:E,onMouseLeave:T,onOpen:A,style:P,text:U,position:j="top",...y},$)=>{const[Y,H]=m.useState(!1),[V,J]=m.useState(),[L,M]=m.useState(),ne=!s,de=!u,Ae=F=>{window.clearTimeout(V),window.clearTimeout(L);const se=window.setTimeout(()=>{H(!0),A==null||A(F)},f);J(se)},$e=F=>{F.persist(),nf(F)?x==null||x(F):of(F)&&(E==null||E(F)),Ae(F)},De=F=>{window.clearTimeout(V),window.clearTimeout(L);const se=window.setTimeout(()=>{H(!1),v==null||v(F)},g);M(se)},W=F=>{F.persist(),nf(F)?h==null||h(F):of(F)&&(T==null||T(F)),De(F)},ce=ne?W:void 0,xe=ne?$e:void 0,K=de?$e:void 0,ue=de?W:void 0,ge=ne?0:void 0;return B.createElement(h0,{"data-testid":"tooltip-wrapper",onBlur:ce,onFocus:xe,onMouseEnter:K,onMouseLeave:ue,tabIndex:ge},B.createElement(p0,{className:n,"data-testid":"tooltip",position:j,ref:$,show:Y,style:P,...y},U),l)});m0.displayName="Tooltip";const js=ee(Os)`
  white-space: nowrap;
`,Zf=me`
  :focus {
    outline: none;
  }

  ${({$disabled:n})=>n?"cursor: default;":me`
          cursor: pointer;

          :focus {
            ${js} {
              background: ${({theme:l})=>l.hoverBackground};
              color: ${({theme:l})=>l.materialTextInvert};
              outline: 2px dotted ${({theme:l})=>l.focusSecondary};
            }
          }
        `}
`,g0=ee.ul`
  position: relative;
  isolation: isolate;

  ${({isRootLevel:n})=>n&&me`
      &:before {
        content: '';
        position: absolute;
        top: 20px;
        bottom: 0;
        left: 5.5px;
        width: 1px;
        border-left: 2px dashed ${({theme:l})=>l.borderDark};
      }
    `}

  ul {
    padding-left: 19.5px;
  }

  li {
    position: relative;

    &:before {
      content: '';
      position: absolute;
      top: 17.5px;
      left: 5.5px;
      width: 22px;
      border-top: 2px dashed ${({theme:n})=>n.borderDark};
      font-size: 12px;
    }
  }
`,v0=ee.li`
  position: relative;
  padding-left: ${({hasItems:n})=>n?"0":"13px"};

  ${({isRootLevel:n})=>n?me`
          &:last-child {
            &:after {
              content: '';
              position: absolute;
              top: 19.5px;
              left: 1px;
              bottom: 0;
              width: 10px;
              background: ${({theme:l})=>l.material};
            }
          }
        `:me`
          &:last-child {
            &:after {
              content: '';
              position: absolute;
              z-index: 1;
              top: 19.5px;
              bottom: 0;
              left: 1.5px;
              width: 10px;
              background: ${({theme:l})=>l.material};
            }
          }
        `}

  & > details > ul {
    &:after {
      content: '';
      position: absolute;
      top: -18px;
      bottom: 0;
      left: 25px;
      border-left: 2px dashed ${({theme:n})=>n.borderDark};
    }
  }
`,y0=ee.details`
  position: relative;
  z-index: 2;

  &[open] > summary:before {
    content: '-';
  }
`,x0=ee.summary`
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  color: ${({theme:n})=>n.materialText};
  user-select: none;
  padding-left: 18px;
  ${Zf};

  &::-webkit-details-marker {
    display: none;
  }

  &:before {
    content: '+';
    position: absolute;
    left: 0;
    display: block;
    width: 8px;
    height: 9px;
    border: 2px solid #808080;
    padding-left: 1px;
    background-color: #fff;
    line-height: 8px;
    text-align: center;
  }
`,lf=ee(Ms)`
  position: relative;
  z-index: 1;
  background: none;
  border: 0;
  font-family: inherit;
  padding-top: 8px;
  padding-bottom: 8px;
  margin: 0;
  ${Zf};
`,w0=ee.span`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-right: 6px;
`;function af(n,l){return n.includes(l)?n.filter(s=>s!==l):[...n,l]}function sf(n){n.preventDefault()}function Jf({className:n,disabled:l,expanded:s,innerRef:u,level:f,select:g,selected:h,style:v,tree:x=[]}){const E=f===0,T=m.useCallback(A=>{var P,U;const j=!!(A.items&&A.items.length>0),y=s.includes(A.id),$=(P=l||A.disabled)!==null&&P!==void 0?P:!1,Y=$?sf:L=>g(L,A),H=$?sf:L=>g(L,A),V=h===A.id,J=B.createElement(w0,{"aria-hidden":!0},A.icon);return B.createElement(v0,{key:A.label,isRootLevel:E,role:"treeitem","aria-expanded":y,"aria-selected":V,hasItems:j},j?B.createElement(y0,{open:y},B.createElement(x0,{onClick:Y,$disabled:$},B.createElement(lf,{$disabled:$},J,B.createElement(js,null,A.label))),y&&B.createElement(Jf,{className:n,disabled:$,expanded:s,level:f+1,select:g,selected:h,style:v,tree:(U=A.items)!==null&&U!==void 0?U:[]})):B.createElement(lf,{as:"button",$disabled:$,onClick:H},J,B.createElement(js,null,A.label)))},[n,l,s,E,f,g,h,v]);return B.createElement(g0,{className:E?n:void 0,style:E?v:void 0,ref:E?u:void 0,role:E?"tree":"group",isRootLevel:E},x.map(T))}function b0({className:n,defaultExpanded:l=[],defaultSelected:s,disabled:u=!1,expanded:f,onNodeSelect:g,onNodeToggle:h,selected:v,style:x,tree:E=[]},T){const[A,P]=Wr({defaultValue:l,onChange:h,onChangePropName:"onNodeToggle",value:f,valuePropName:"expanded"}),[U,j]=Wr({defaultValue:s,onChange:g,onChangePropName:"onNodeSelect",value:v,valuePropName:"selected"}),y=m.useCallback((H,V)=>{if(h){const J=af(A,V);h(H,J)}P(J=>af(J,V))},[A,h,P]),$=m.useCallback((H,V)=>{j(V),g&&g(H,V)},[g,j]),Y=m.useCallback((H,V)=>{H.preventDefault(),$(H,V.id),V.items&&V.items.length&&y(H,V.id)},[$,y]);return B.createElement(Jf,{className:n,disabled:u,expanded:A,level:0,innerRef:T,select:Y,selected:U,style:x,tree:E})}const k0=m.forwardRef(b0);k0.displayName="TreeView";const S0="/assets/ms_sans_serif-Du8rjN1q.woff2",C0="/assets/ms_sans_serif_bold-D5dpRRHG.woff2";var ls,cf;function A0(){if(cf)return ls;cf=1;var n={name:"original",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#848584",borderDarkest:"#0a0a0a",borderLight:"#dfdfdf",borderLightest:"#fefefe",canvas:"#ffffff",canvasText:"#0a0a0a",canvasTextDisabled:"#848584",canvasTextDisabledShadow:"#fefefe",canvasTextInvert:"#fefefe",checkmark:"#0a0a0a",checkmarkDisabled:"#848584",desktopBackground:"#008080",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#fefe03",headerBackground:"#060084",headerNotActiveBackground:"#7f787f",headerNotActiveText:"#c6c6c6",headerText:"#fefefe",hoverBackground:"#060084",material:"#c6c6c6",materialDark:"#9a9e9c",materialText:"#0a0a0a",materialTextDisabled:"#848584",materialTextDisabledShadow:"#fefefe",materialTextInvert:"#fefefe",progress:"#060084",tooltip:"#fefbcc"};return ls=n,ls}var E0=A0();const _0=Ut(E0);var as,uf;function T0(){if(uf)return as;uf=1;var n={name:"tokyoDark",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#1f2223",borderDarkest:"#070809",borderLight:"#5e696a",borderLightest:"#93a0a1",canvas:"#2f3435",canvasText:"#F4F4ED",canvasTextDisabled:"#1f2223",canvasTextDisabledShadow:"#93a0a1",canvasTextInvert:"#ffffff",checkmark:"#F4F4ED",checkmarkDisabled:"#1f2223",desktopBackground:"#181a1b",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#20FC8F",headerBackground:"#1f2223",headerNotActiveBackground:"#5e696a",headerNotActiveText:"#F4F4ED",headerText:"#F4F4ED",hoverBackground:"#F61067",material:"#465051",materialDark:"#1f2223",materialText:"#F4F4ED",materialTextDisabled:"#1f2223",materialTextDisabledShadow:"#93a0a1",materialTextInvert:"#ffffff",progress:"#F61067",tooltip:"#fefbcc"};return as=n,as}var j0=T0();const P0=Ut(j0);var ss,df;function R0(){if(df)return ss;df=1;var n={name:"vaporTeal",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#00706f",borderDarkest:"#000000",borderLight:"#2fcecd",borderLightest:"#58ffff",canvas:"#98DFEA",canvasText:"#000000",canvasTextDisabled:"#00706f",canvasTextDisabledShadow:"#58ffff",canvasTextInvert:"#000000",checkmark:"#000000",checkmarkDisabled:"#00706f",desktopBackground:"#008080",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#FCF6BD",headerBackground:"#246A73",headerNotActiveBackground:"#2fcecd",headerNotActiveText:"#246A73",headerText:"#58ffff",hoverBackground:"#FF99C8",material:"#01a8a8",materialDark:"#246A73",materialText:"#000000",materialTextDisabled:"#00706f",materialTextDisabledShadow:"#58ffff",materialTextInvert:"#000000",progress:"#FF99C8",tooltip:"#fefbcc"};return ss=n,ss}var $0=R0();const D0=Ut($0);var cs,ff;function z0(){if(ff)return cs;ff=1;var n={name:"rainyDay",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#3d5367",borderDarkest:"#16233b",borderLight:"#91abc2",borderLightest:"#b7cee5",canvas:"#ffffff",canvasText:"#050608",canvasTextDisabled:"#3d5367",canvasTextDisabledShadow:"#b7cee5",canvasTextInvert:"#ffffff",checkmark:"#050608",checkmarkDisabled:"#3d5367",desktopBackground:"#000000",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#fefe03",headerBackground:"#4b6480",headerNotActiveBackground:"#7f7f81",headerNotActiveText:"#ced0d9",headerText:"#ffffff",hoverBackground:"#4b6480",material:"#7a99b3",materialDark:"#9a9e9c",materialText:"#050608",materialTextDisabled:"#3d5367",materialTextDisabledShadow:"#b7cee5",materialTextInvert:"#ffffff",progress:"#4b6480",tooltip:"#fefbcc"};return cs=n,cs}var I0=z0();const L0=Ut(I0);var us,pf;function M0(){if(pf)return us;pf=1;var n={name:"modernDark",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#121317",borderDarkest:"#000000",borderLight:"#31323c",borderLightest:"#4b4d58",canvas:"#4b4d58",canvasText:"#000000",canvasTextDisabled:"#4b4d58",canvasTextDisabledShadow:"#4b4d58",canvasTextInvert:"#202127",checkmark:"#000000",checkmarkDisabled:"#121317",desktopBackground:"#000000",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#fefe03",headerBackground:"#4b4d58",headerNotActiveBackground:"transparent",headerNotActiveText:"#4b4d58",headerText:"#202127",hoverBackground:"#f88702",material:"#202127",materialDark:"#9a9e9c",materialText:"#f88702",materialTextDisabled:"#4b4d58",materialTextDisabledShadow:"#121317",materialTextInvert:"#202127",progress:"#f88702",tooltip:"#fefbcc"};return us=n,us}var O0=M0();const N0=Ut(O0);var ds,hf;function B0(){if(hf)return ds;hf=1;var n={name:"matrix",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#282828",borderDarkest:"#000000",borderLight:"#656565",borderLightest:"#a7a7a7",canvas:"#c0c0c0",canvasText:"#000000",canvasTextDisabled:"#282828",canvasTextDisabledShadow:"#ff0000",canvasTextInvert:"#ffffff",checkmark:"#000000",checkmarkDisabled:"#282828",desktopBackground:"#000000",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#35FF69",headerBackground:"#000000",headerNotActiveBackground:"#7f7f7f",headerNotActiveText:"#535353",headerText:"#a7a7a7",hoverBackground:"#000000",material:"#535353",materialDark:"#282828",materialText:"#35FF69",materialTextDisabled:"#282828",materialTextDisabledShadow:"#a7a7a7",materialTextInvert:"#ffffff",progress:"#000000",tooltip:"#fefbcc"};return ds=n,ds}var F0=B0();const W0=Ut(F0);var fs,mf;function H0(){if(mf)return fs;mf=1;var n={name:"lilac",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#5f549b",borderDarkest:"#1c1449",borderLight:"#bcb4e9",borderLightest:"#d3ccf4",canvas:"#ffffff",canvasText:"#050608",canvasTextDisabled:"#5f549b",canvasTextDisabledShadow:"#ffffff",canvasTextInvert:"#ffffff",checkmark:"#050608",checkmarkDisabled:"#5f549b",desktopBackground:"#000000",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#fefe03",headerBackground:"#5e4dba",headerNotActiveBackground:"#7f7f81",headerNotActiveText:"#ced0cf",headerText:"#ffffff",hoverBackground:"#5e4dba",material:"#b1a7df",materialDark:"#9a9e9c",materialText:"#050608",materialTextDisabled:"#5f549b",materialTextDisabledShadow:"#ffffff",materialTextInvert:"#ffffff",progress:"#5e4dba",tooltip:"#fefbcc"};return fs=n,fs}var U0=H0();const V0=Ut(U0);var ps,gf;function Q0(){if(gf)return ps;gf=1;var n={name:"rose",anchor:"#1034a6",anchorVisited:"#440381",borderDark:"#8a5b68",borderDarkest:"#26030b",borderLight:"#e5bec8",borderLightest:"#f1d4dc",canvas:"#ffffff",canvasText:"#050608",canvasTextDisabled:"#8a5b68",canvasTextDisabledShadow:"#f1d4dc",canvasTextInvert:"#ffffff",checkmark:"#050608",checkmarkDisabled:"#8a5b68",desktopBackground:"#808080",flatDark:"#9e9e9e",flatLight:"#d8d8d8",focusSecondary:"#fefe03",headerBackground:"#ab5a71",headerNotActiveBackground:"#a19fa5",headerNotActiveText:"#615f68",headerText:"#ffffff",hoverBackground:"#ab5a71",material:"#d6adb8",materialDark:"#9a9e9c",materialText:"#050608",materialTextDisabled:"#8a5b68",materialTextDisabledShadow:"#f1d4dc",materialTextInvert:"#ffffff",progress:"#ab5a71",tooltip:"#fefbcc"};return ps=n,ps}var G0=Q0();const K0=Ut(G0);var hs,vf;function Y0(){if(vf)return hs;vf=1;var n={name:"hotdogStand",anchor:"rgb(255, 255, 0)",anchorVisited:"rgb(255, 255, 0)",borderDark:"rgb(0, 0, 0)",borderDarkest:"rgb(0, 0, 0)",borderLight:"rgb(0, 0, 0)",borderLightest:"rgb(0, 0, 0)",canvas:"rgb(255, 0, 0)",canvasText:"rgb(255, 255, 255)",canvasTextDisabled:"rgb(0, 0, 0)",canvasTextDisabledShadow:"rgb(0, 0, 0)",canvasTextInvert:"rgb(255, 255, 255)",checkmark:"rgb(255, 255, 255)",checkmarkDisabled:"rgb(128, 128, 128)",desktopBackground:"rgb(255, 255, 0)",flatDark:"rgb(0, 0, 0)",flatLight:"rgb(0, 0, 0)",focusSecondary:"rgb(0, 0, 0)",headerBackground:"linear-gradient(to right, rgb(0, 0, 0), rgb(0, 0, 0))",headerNotActiveBackground:"linear-gradient(to right, rgb(255, 0, 0), rgb(255, 0, 0))",headerNotActiveText:"rgb(255, 255, 255)",headerText:"rgb(255, 255, 255)",hoverBackground:"rgb(0, 0, 0)",material:"rgb(255, 0, 0)",materialDark:"rgb(255, 0, 0)",materialText:"rgb(255, 255, 255)",materialTextDisabled:"rgb(0, 0, 0)",materialTextDisabledShadow:"rgb(0, 0, 0)",materialTextInvert:"rgb(255, 255, 255)",progress:"rgb(0, 0, 0)",tooltip:"rgb(255, 255, 225)"};return hs=n,hs}var q0=Y0();const X0=Ut(q0);async function He(n,l){const s={...l==null?void 0:l.headers};l!=null&&l.method&&["POST","PUT","PATCH"].includes(l.method.toUpperCase())&&(s["Content-Type"]="application/json");const u=await fetch(n,{...l,headers:s}),f=await u.text();let g;if(f)try{g=JSON.parse(f)}catch{if(!u.ok)throw new Error(f||`HTTP ${u.status}`)}if(!u.ok){const h=typeof g=="object"&&g!==null&&"error"in g?String(g.error):f||`HTTP ${u.status}`;throw new Error(h)}return g}const St=encodeURIComponent,Le={info:()=>He("/api/info"),state:()=>He("/api/state"),updateAgents:n=>He("/api/agents",{method:"PUT",body:JSON.stringify(n)}),agentConfigs:()=>He("/api/agents/configs"),createService:n=>He("/api/services",{method:"POST",body:JSON.stringify(n)}),updateService:(n,l)=>He(`/api/services/${St(n)}`,{method:"PUT",body:JSON.stringify(l)}),deleteService:(n,l=!0)=>He(`/api/services/${St(n)}?deregister=${l?"true":"false"}`,{method:"DELETE"}),startService:(n,l=!0)=>He(`/api/services/${St(n)}/start`,{method:"POST",body:JSON.stringify({build:l})}),stopService:(n,l=!0)=>He(`/api/services/${St(n)}/stop`,{method:"POST",body:JSON.stringify({deregister:l})}),restartService:n=>He(`/api/services/${St(n)}/restart`,{method:"POST",body:"{}"}),rebuildService:n=>He(`/api/services/${St(n)}/rebuild`,{method:"POST",body:"{}"}),registerService:n=>He(`/api/services/${St(n)}/register`,{method:"POST",body:"{}"}),deregisterService:n=>He(`/api/services/${St(n)}/deregister`,{method:"POST",body:"{}"}),healthService:n=>He(`/api/services/${St(n)}/health`),getTools:n=>He(`/api/services/${St(n)}/tools`),syncTools:n=>He(`/api/services/${St(n)}/tools/sync`,{method:"POST",body:"{}"}),updateTools:(n,l)=>He(`/api/services/${St(n)}/tools`,{method:"PUT",body:JSON.stringify(l)}),applyTools:n=>He(`/api/services/${St(n)}/tools/apply`,{method:"POST",body:"{}"}),getLogs:(n,l=200)=>He(`/api/services/${St(n)}/logs?lines=${l}`),startAll:()=>He("/api/actions/start-all",{method:"POST",body:"{}"}),stopAll:(n=!0)=>He("/api/actions/stop-all",{method:"POST",body:JSON.stringify({deregister:n})}),restartAll:()=>He("/api/actions/restart-all",{method:"POST",body:"{}"}),registerAll:()=>He("/api/actions/register-all",{method:"POST",body:"{}"}),validate:()=>He("/api/validate"),doctor:(n="",l=!1)=>He("/api/doctor",{method:"POST",body:JSON.stringify({service:n,fix:l})}),importPreview:()=>He("/api/import/preview"),importApply:(n,l=!1)=>He("/api/import",{method:"POST",body:JSON.stringify({items:n,overwrite:l})}),getRawConfig:()=>He("/api/config/raw"),putRawConfig:n=>He("/api/config/raw",{method:"PUT",body:JSON.stringify({content:n})}),searchRegistry:n=>He(`/api/registry/search?q=${St(n)}`)};function ep(){return{name:"",type:"",command:"",args:[],env:{},port:0,extra_ports:[],build_command:"",no_build_on_start:!1,deps:[],description:"",health_url:"",dashboard_url:"",mcp_url:"",headers:{},log:"",active_tier:"",code_mode:!1,no_code_mode:!1,tools_config_path:"",path:"",installer:""}}function Z0(n){return{name:n.name||"",type:n.type||"",command:n.command||"",args:n.args?[...n.args]:[],env:n.env?{...n.env}:{},port:n.port||0,extra_ports:n.extra_ports?[...n.extra_ports]:[],build_command:n.effective_build||n.build_command||n.build_cmd||"",no_build_on_start:!!n.no_build_on_start,deps:n.deps?[...n.deps]:[],description:n.description||"",health_url:n.health_url||"",dashboard_url:n.dashboard_url||"",mcp_url:n.mcp_url||"",headers:n.headers?{...n.headers}:{},log:n.log||"",active_tier:n.active_tier||"",code_mode:!!n.code_mode,no_code_mode:!!n.no_code_mode,tools_config_path:n.tools_config_path||"",path:n.path||"",installer:n.installer||""}}function qi(n){return n?Object.entries(n).map(([l,s])=>({key:l,value:s})):[]}function yf(n){const l={};for(const s of n){const u=s.key.trim();u&&(l[u]=s.value)}return l}const J0=[{value:"",label:"Auto-detect"},{value:"http",label:"http (Streamable HTTP)"},{value:"stdio",label:"stdio (Agent-launched)"},{value:"filesystem",label:"filesystem (Embedded FS)"}],ey=[{value:"",label:"(default)"},{value:"core",label:"core"},{value:"extended",label:"extended"},{value:"complete",label:"complete"}],ty=[{value:"",label:"Auto-detect"},{value:"ast-mcp",label:"ast-mcp (delegated)"},{value:"native",label:"native (mcp-local)"}],tp=({initial:n,isNew:l=!1,saving:s=!1,submitLabel:u="Save Configuration",onSubmit:f,onCancel:g,extraButtons:h})=>{const[v,x]=m.useState(n),[E,T]=m.useState(n.port?String(n.port):""),[A,P]=m.useState((n.extra_ports||[]).join(", ")),[U,j]=m.useState(n.args||[]),[y,$]=m.useState(n.deps||[]),[Y,H]=m.useState(()=>qi(n.env)),[V,J]=m.useState(()=>qi(n.headers)),[L,M]=m.useState(!1),[ne,de]=m.useState(0),[Ae,$e]=m.useState("");m.useEffect(()=>{x(n),T(n.port?String(n.port):""),P((n.extra_ports||[]).join(", ")),j(n.args||[]),$(n.deps||[]),H(qi(n.env)),J(qi(n.headers)),$e("")},[n]);const De=async W=>{W.preventDefault(),$e("");const ce=E.trim()===""?0:Number(E.trim());if(Number.isNaN(ce)||ce<0||ce>65535){$e("Port must be a number between 0 and 65535.");return}const xe=[];if(A.trim()!=="")for(const ue of A.split(",")){const ge=Number(ue.trim());if(Number.isNaN(ge)||ge<=0||ge>65535){$e(`Invalid extra port: "${ue.trim()}"`);return}xe.push(ge)}const K={...v,name:v.name.trim(),port:ce,extra_ports:xe,args:U.filter(ue=>ue!==""),deps:y.map(ue=>ue.trim()).filter(Boolean),env:yf(Y),headers:yf(V)};if(!K.name){$e("Service name is required.");return}await f(K)};return c.jsxs("form",{onSubmit:De,style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsxs(ll,{value:ne,onChange:W=>de(W),children:[c.jsx(ar,{value:0,children:"General & Command"}),c.jsx(ar,{value:1,children:"Network & URLs"}),c.jsxs(ar,{value:2,children:["Env & Headers (",Y.length+V.length,")"]}),c.jsx(ar,{value:3,children:"AST & Advanced"})]}),c.jsxs(il,{style:{padding:12},children:[ne===0&&c.jsxs("div",{style:{display:"grid",gridTemplateColumns:"repeat(auto-fit, minmax(230px, 1fr))",gap:10},children:[c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Service Name *"}),c.jsx(Ke,{value:v.name,onChange:W=>x({...v,name:W.target.value}),placeholder:"e.g. ast-context-cache",fullWidth:!0})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Transport / Type"}),c.jsx(Nn,{options:J0,value:v.type,onChange:W=>x({...v,type:W.value}),width:"100%"})]}),c.jsxs("div",{style:{gridColumn:"1 / -1"},children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Description"}),c.jsx(Ke,{value:v.description,onChange:W=>x({...v,description:W.target.value}),placeholder:"Optional human-readable summary",fullWidth:!0})]}),v.type==="filesystem"?c.jsxs("div",{style:{gridColumn:"1 / -1"},children:[c.jsxs("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:["Root Directory Path * (supports ~/ and $","{HOME}",")"]}),c.jsx(Ke,{value:v.path,onChange:W=>x({...v,path:W.target.value}),placeholder:"~/projects",fullWidth:!0})]}):c.jsxs(c.Fragment,{children:[c.jsxs("div",{style:{gridColumn:"1 / -1"},children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Command (binary path or npx/uvx; blank for remote-only)"}),c.jsx(Ke,{value:v.command,onChange:W=>x({...v,command:W.target.value}),placeholder:"${HOME}/git/ast-context-cache/ast-mcp or npx",fullWidth:!0})]}),c.jsx("div",{style:{gridColumn:"1 / -1"},children:c.jsx(Ct,{label:"Arguments (args)",children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:6},children:[U.length===0&&c.jsx("div",{style:{fontSize:12,opacity:.75},children:"No arguments configured."}),U.map((W,ce)=>c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center"},children:[c.jsx(Ke,{value:W,onChange:xe=>{const K=[...U];K[ce]=xe.target.value,j(K)},placeholder:`arg[${ce}]`,fullWidth:!0}),c.jsx(he,{type:"button",size:"sm",onClick:()=>j(U.filter((xe,K)=>K!==ce)),children:"Remove"})]},ce)),c.jsx("div",{children:c.jsx(he,{type:"button",size:"sm",onClick:()=>j([...U,""]),children:"+ Add Arg"})})]})})})]}),c.jsxs("div",{style:{gridColumn:"1 / -1",display:"flex",gap:10,alignItems:"flex-end",flexWrap:"wrap"},children:[c.jsxs("div",{style:{flex:1,minWidth:220},children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Build Command"}),c.jsx(Ke,{value:v.build_command,onChange:W=>x({...v,build_command:W.target.value}),placeholder:"make build",fullWidth:!0})]}),c.jsx(ut,{checked:v.no_build_on_start,onChange:W=>x({...v,no_build_on_start:W.target.checked}),label:"Skip build on start (no_build_on_start)"})]})]}),ne===1&&c.jsxs("div",{style:{display:"grid",gridTemplateColumns:"repeat(auto-fit, minmax(230px, 1fr))",gap:10},children:[c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Primary Port (0 for stdio)"}),c.jsx(Ke,{value:E,onChange:W=>T(W.target.value),placeholder:"7821",fullWidth:!0})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Extra Ports (comma-separated)"}),c.jsx(Ke,{value:A,onChange:W=>P(W.target.value),placeholder:"7830",fullWidth:!0})]}),c.jsxs("div",{style:{gridColumn:"1 / -1"},children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"MCP Endpoint URL (mcp_url)"}),c.jsx(Ke,{value:v.mcp_url,onChange:W=>x({...v,mcp_url:W.target.value}),placeholder:"http://localhost:7821/mcp",fullWidth:!0})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Health URL (health_url)"}),c.jsx(Ke,{value:v.health_url,onChange:W=>x({...v,health_url:W.target.value}),placeholder:"http://localhost:7821/health",fullWidth:!0})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Dashboard URL (dashboard_url)"}),c.jsx(Ke,{value:v.dashboard_url,onChange:W=>x({...v,dashboard_url:W.target.value}),placeholder:"http://localhost:7830",fullWidth:!0})]})]}),ne===2&&c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:12},children:[c.jsx(Ct,{label:"Environment Variables (env)",children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:6},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center"},children:[c.jsx(ut,{checked:L,onChange:W=>M(W.target.checked),label:"Reveal secret values"}),c.jsx(he,{type:"button",size:"sm",onClick:()=>H([...Y,{key:"",value:""}]),children:"+ Add Env Var"})]}),Y.length===0&&c.jsx("div",{style:{fontSize:12,opacity:.75},children:"No environment variables configured."}),Y.map((W,ce)=>{const xe=/token|secret|key|pass|auth/i.test(W.key);return c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center"},children:[c.jsx(Ke,{value:W.key,onChange:K=>{const ue=[...Y];ue[ce]={...W,key:K.target.value},H(ue)},placeholder:"KEY",style:{flex:1}}),c.jsx(Ke,{type:xe&&!L?"password":"text",value:W.value,onChange:K=>{const ue=[...Y];ue[ce]={...W,value:K.target.value},H(ue)},placeholder:"VALUE",style:{flex:2}}),c.jsx(he,{type:"button",size:"sm",onClick:()=>H(Y.filter((K,ue)=>ue!==ce)),children:"×"})]},ce)})]})}),c.jsx(Ct,{label:"HTTP Headers (headers)",children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:6},children:[V.length===0&&c.jsx("div",{style:{fontSize:12,opacity:.75},children:"No custom HTTP headers configured."}),V.map((W,ce)=>c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center"},children:[c.jsx(Ke,{value:W.key,onChange:xe=>{const K=[...V];K[ce]={...W,key:xe.target.value},J(K)},placeholder:"Authorization",style:{flex:1}}),c.jsx(Ke,{type:L?"text":"password",value:W.value,onChange:xe=>{const K=[...V];K[ce]={...W,value:xe.target.value},J(K)},placeholder:"Bearer ...",style:{flex:2}}),c.jsx(he,{type:"button",size:"sm",onClick:()=>J(V.filter((xe,K)=>K!==ce)),children:"×"})]},ce)),c.jsx("div",{children:c.jsx(he,{type:"button",size:"sm",onClick:()=>J([...V,{key:"",value:""}]),children:"+ Add Header"})})]})})]}),ne===3&&c.jsxs("div",{style:{display:"grid",gridTemplateColumns:"repeat(auto-fit, minmax(230px, 1fr))",gap:10},children:[c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Log File Path (blank = ~/.mcp-local/logs/<name>.log)"}),c.jsx(Ke,{value:v.log,onChange:W=>x({...v,log:W.target.value}),placeholder:"~/logs/service.log",fullWidth:!0})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Installer"}),c.jsx(Nn,{options:ty,value:v.installer,onChange:W=>x({...v,installer:W.value}),width:"100%"})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"AST Active Tier (active_tier)"}),c.jsx(Nn,{options:ey,value:v.active_tier,onChange:W=>x({...v,active_tier:W.value}),width:"100%"})]}),c.jsxs("div",{children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Tools Config Path (tools_config_path)"}),c.jsx(Ke,{value:v.tools_config_path,onChange:W=>x({...v,tools_config_path:W.target.value}),placeholder:"~/.astcache/tools.json",fullWidth:!0})]}),c.jsxs("div",{style:{gridColumn:"1 / -1",display:"flex",gap:16,flexWrap:"wrap"},children:[c.jsx(ut,{checked:v.code_mode,onChange:W=>x({...v,code_mode:W.target.checked,no_code_mode:W.target.checked?!1:v.no_code_mode}),label:"Enable Code Mode (AST_MCP_CODE_MODE=true)"}),c.jsx(ut,{checked:v.no_code_mode,onChange:W=>x({...v,no_code_mode:W.target.checked,code_mode:W.target.checked?!1:v.code_mode}),label:"Disable Code Mode (no_code_mode)"})]}),c.jsx("div",{style:{gridColumn:"1 / -1"},children:c.jsx(Ct,{label:"Required File Dependencies (deps)",children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:6},children:[y.length===0&&c.jsx("div",{style:{fontSize:12,opacity:.75},children:"No dependency paths checked before start."}),y.map((W,ce)=>c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center"},children:[c.jsx(Ke,{value:W,onChange:xe=>{const K=[...y];K[ce]=xe.target.value,$(K)},placeholder:"/path/to/required/file",fullWidth:!0}),c.jsx(he,{type:"button",size:"sm",onClick:()=>$(y.filter((xe,K)=>K!==ce)),children:"Remove"})]},ce)),c.jsx("div",{children:c.jsx(he,{type:"button",size:"sm",onClick:()=>$([...y,""]),children:"+ Add Dependency"})})]})})})]})]}),Ae&&c.jsx("div",{style:{color:"#c00000",fontWeight:"bold",fontSize:12},children:Ae}),c.jsxs("div",{style:{display:"flex",justifyContent:"flex-end",gap:8,flexWrap:"wrap"},children:[h,g&&c.jsx(he,{type:"button",onClick:g,disabled:s,children:"Cancel"}),c.jsx(he,{type:"submit",primary:!0,disabled:s,children:s?"Saving...":l?"Create Service":u})]})]})};function be(n){return{...ep(),...n,args:n.args?[...n.args]:[],env:n.env?{...n.env}:{},headers:n.headers?{...n.headers}:{},deps:n.deps?[...n.deps]:[],extra_ports:n.extra_ports?[...n.extra_ports]:[]}}const ry=["All Categories","Zapier (9,000+ Apps)","Coding & Git","Browser & Automation","Search & Web","Databases & Storage","Productivity & Docs","Cloud & DevOps","AI & Memory","Templates"],ny=[{id:"all",label:"All Sources"},{id:"popular",label:"★ Top Popular"},{id:"github.com/mcp",label:"github.com/mcp"},{id:"mcpservers.org",label:"mcpservers.org"},{id:"zapier.com/mcp",label:"zapier.com/mcp"},{id:"template",label:"Starter Templates"}],ms=[{id:"zapier-mcp-all",title:"Zapier MCP — All 9,000+ Apps (Unified Endpoint)",category:"Zapier (9,000+ Apps)",source:"zapier.com/mcp",source_url:"https://zapier.com/mcp",popular:!0,description:"Connect AI agents to 9,000+ apps and 30,000+ actions via Zapier Streamable HTTP MCP (Slack, Gmail, HubSpot, Notion, Jira, Google Sheets, Salesforce, etc.).",tags:["zapier","automation","slack","gmail","google sheets","calendar","hubspot","salesforce","jira","trello","airtable","remote","http"],setup_hint:"Create a server at https://mcp.zapier.com and paste your Connection Token into the Authorization Bearer header (or replace mcp_url with your server-specific URL).",service:be({name:"zapier-mcp",type:"http",mcp_url:"https://mcp.zapier.com/api/v1/connect",headers:{Authorization:"Bearer YOUR_ZAPIER_CONNECTION_TOKEN"},dashboard_url:"https://mcp.zapier.com",description:"Zapier MCP gateway (9,000+ apps via Streamable HTTP)"})},{id:"zapier-mcp-google-workspace",title:"Zapier MCP — Gmail, Calendar, Sheets & Drive Bundle",category:"Zapier (9,000+ Apps)",source:"zapier.com/mcp",source_url:"https://zapier.com/mcp",popular:!0,description:"Send/draft Gmail emails, schedule Google Calendar events, append/query Google Sheets rows, and manage Google Drive files via Zapier MCP.",tags:["zapier","gmail","google calendar","google sheets","google docs","google drive","email"],setup_hint:"Enable Gmail, Google Calendar, Google Sheets, and Google Drive actions in your Zapier MCP server at https://mcp.zapier.com and set your Bearer token.",service:be({name:"zapier-google-workspace",type:"http",mcp_url:"https://mcp.zapier.com/api/v1/connect",headers:{Authorization:"Bearer YOUR_ZAPIER_CONNECTION_TOKEN"},dashboard_url:"https://mcp.zapier.com",description:"Zapier MCP bundle for Gmail, Google Calendar, Sheets & Drive"})},{id:"zapier-mcp-crm",title:"Zapier MCP — CRM & Sales Bundle (HubSpot, Salesforce, Pipedrive)",category:"Zapier (9,000+ Apps)",source:"zapier.com/mcp",source_url:"https://zapier.com/mcp",description:"Look up contacts, update deals, log notes, and trigger CRM workflows across HubSpot, Salesforce, Pipedrive, and Zoho via Zapier MCP.",tags:["zapier","crm","hubspot","salesforce","pipedrive","zoho","sales","leads"],setup_hint:"Create a dedicated CRM tool bundle at https://mcp.zapier.com and paste its Connection Token in the Authorization header.",service:be({name:"zapier-crm",type:"http",mcp_url:"https://mcp.zapier.com/api/v1/connect",headers:{Authorization:"Bearer YOUR_ZAPIER_CONNECTION_TOKEN"},dashboard_url:"https://mcp.zapier.com",description:"Zapier MCP CRM bundle (HubSpot, Salesforce, Pipedrive)"})},{id:"zapier-mcp-pm",title:"Zapier MCP — Project Management (Jira, Asana, Trello, ClickUp, Monday)",category:"Zapier (9,000+ Apps)",source:"zapier.com/mcp",source_url:"https://zapier.com/mcp",description:"Create issues, update boards, assign tasks, and sync sprints across Jira, Trello, Asana, ClickUp, and Monday.com via Zapier MCP.",tags:["zapier","jira","asana","trello","clickup","monday","project management","tasks"],setup_hint:"Add your project management tool actions in https://mcp.zapier.com and paste your Connection Token.",service:be({name:"zapier-pm",type:"http",mcp_url:"https://mcp.zapier.com/api/v1/connect",headers:{Authorization:"Bearer YOUR_ZAPIER_CONNECTION_TOKEN"},dashboard_url:"https://mcp.zapier.com",description:"Zapier MCP Project Management bundle (Jira, Asana, Trello, ClickUp)"})},{id:"zapier-mcp-comms",title:"Zapier MCP — Messaging & Comms (Slack, Teams, Discord, Twilio)",category:"Zapier (9,000+ Apps)",source:"zapier.com/mcp",source_url:"https://zapier.com/mcp",description:"Post channel messages, send DMs, search threads, and send SMS notifications across Slack, Microsoft Teams, Discord, and Twilio.",tags:["zapier","slack","microsoft teams","discord","twilio","chat","messaging"],setup_hint:"Configure Slack / Teams / Discord actions in https://mcp.zapier.com and paste your Connection Token.",service:be({name:"zapier-comms",type:"http",mcp_url:"https://mcp.zapier.com/api/v1/connect",headers:{Authorization:"Bearer YOUR_ZAPIER_CONNECTION_TOKEN"},dashboard_url:"https://mcp.zapier.com",description:"Zapier MCP Messaging bundle (Slack, Teams, Discord)"})},{id:"github-stdio",title:"GitHub MCP Server (Local stdio)",category:"Coding & Git",source:"github.com/mcp",source_url:"https://github.com/github/github-mcp-server",popular:!0,description:"Official GitHub MCP integration: repositories, pull requests, issues, code search, GitHub Actions workflows, and security alerts.",tags:["github","git","pull requests","issues","actions","code search","official"],setup_hint:"Set GITHUB_PERSONAL_ACCESS_TOKEN in Environment Variables with repo/read permissions.",service:be({name:"github",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-github"],env:{GITHUB_PERSONAL_ACCESS_TOKEN:""},description:"GitHub MCP server (PRs, issues, repos, code search)"})},{id:"github-remote",title:"GitHub Copilot Remote MCP Server (Streamable HTTP)",category:"Coding & Git",source:"github.com/mcp",source_url:"https://github.com/github/github-mcp-server",popular:!0,description:"GitHub hosted remote MCP server at api.githubcopilot.com/mcp — zero local install required.",tags:["github","copilot","remote","http","pull requests","issues"],setup_hint:"Paste your GitHub Personal Access Token after Bearer in the Authorization header.",service:be({name:"github-remote",type:"http",mcp_url:"https://api.githubcopilot.com/mcp/",headers:{Authorization:"Bearer YOUR_GITHUB_PAT"},description:"Hosted GitHub Remote MCP server"})},{id:"context7-stdio",title:"Context7 — Up-to-date Library Docs (Upstash)",category:"Coding & Git",source:"github.com/mcp",source_url:"https://github.com/upstash/context7",popular:!0,description:"Fetches up-to-date, version-specific documentation and code examples for thousands of libraries directly into your prompt.",tags:["context7","upstash","docs","documentation","libraries","sdk","code examples"],setup_hint:"Works out of the box via npx; optional CONTEXT7_API_KEY for higher rate limits.",service:be({name:"context7",type:"stdio",command:"npx",args:["-y","@upstash/context7-mcp@latest"],description:"Context7 live library documentation & code examples"})},{id:"context7-remote",title:"Context7 Remote MCP (Streamable HTTP)",category:"Coding & Git",source:"github.com/mcp",source_url:"https://context7.com",description:"Hosted Context7 Streamable HTTP MCP endpoint for version-specific library documentation.",tags:["context7","remote","http","docs","upstash"],service:be({name:"context7-remote",type:"http",mcp_url:"https://mcp.context7.com/mcp",description:"Context7 hosted remote MCP server"})},{id:"ast-context-cache",title:"ast-context-cache (ast-mcp Local Code Index & Search)",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/coma-toast/ast-context-cache",popular:!0,description:"Token-efficient AST code search, call graphs, symbol outlines, virtual context, and subagent handoffs with tool-tier management.",tags:["ast","ast-mcp","code search","symbols","tree-sitter","local"],setup_hint:"Supports active_tier (core/extended/complete) and ~/.astcache/tools.json overrides.",service:be({name:"ast-context-cache",type:"http",command:"${HOME}/.local/bin/ast-mcp",port:7821,mcp_url:"http://localhost:7821/mcp",health_url:"http://localhost:7821/health",dashboard_url:"http://localhost:7821",active_tier:"core",installer:"ast-mcp",description:"AST Context Cache MCP server"})},{id:"git-mcp",title:"Git MCP Server (Official Local Repo Tools)",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Read, search, diff, log, branch, and manipulate local Git repositories via uvx mcp-server-git.",tags:["git","commits","diff","branches","repo","official"],setup_hint:"Optionally pass --repository /path/to/repo in Args.",service:be({name:"git",type:"stdio",command:"uvx",args:["mcp-server-git"],description:"Official Git repository inspection & operations MCP"})},{id:"gitlab-mcp",title:"GitLab MCP Server",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Manage GitLab projects, merge requests, issues, pipelines, and repository files.",tags:["gitlab","git","merge requests","ci/cd","issues"],setup_hint:"Set GITLAB_PERSONAL_ACCESS_TOKEN (and optional GITLAB_API_URL for self-hosted GitLab).",service:be({name:"gitlab",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-gitlab"],env:{GITLAB_PERSONAL_ACCESS_TOKEN:"",GITLAB_API_URL:"https://gitlab.com/api/v4"},description:"GitLab API integration (MRs, issues, repos)"})},{id:"serena-mcp",title:"Serena — Semantic Coding Agent Toolkit (LSP)",category:"Coding & Git",source:"github.com/mcp",source_url:"https://github.com/oraios/serena",popular:!0,description:"Language-server-powered semantic code retrieval and editing tools (find symbol, references, replace symbol body).",tags:["serena","lsp","refactoring","coding agent","symbols","python","typescript","go"],service:be({name:"serena",type:"stdio",command:"uvx",args:["--from","git+https://github.com/oraios/serena","serena-mcp-server"],description:"Serena LSP-powered semantic code search & editing MCP"})},{id:"desktop-commander",title:"Desktop Commander — Terminal & File Diff Editing",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/wonderwhy-er/DesktopCommanderMCP",popular:!0,description:"Execute terminal commands, manage long-running processes, and perform surgical search/replace file edits.",tags:["desktop-commander","terminal","shell","files","diff","processes"],service:be({name:"desktop-commander",type:"stdio",command:"npx",args:["-y","@wonderwhy-er/desktop-commander@latest"],description:"Terminal process control and diff-based file editing"})},{id:"e2b-sandbox",title:"E2B Code Interpreter Sandbox",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/e2b-dev/mcp-server",description:"Run Python and JavaScript code securely in isolated cloud sandboxes via E2B.",tags:["e2b","sandbox","code interpreter","python","execution"],setup_hint:"Requires E2B_API_KEY in Environment Variables.",service:be({name:"e2b",type:"stdio",command:"npx",args:["-y","@e2b/mcp-server"],env:{E2B_API_KEY:""},description:"E2B secure cloud code execution sandbox"})},{id:"filesystem-embedded",title:"Filesystem (mcp-local Embedded HTTP Server)",category:"Coding & Git",source:"template",popular:!0,description:"Built-in mcp-local HTTP filesystem server exposing a directory tree to all enabled agents.",tags:["filesystem","files","local","directory","embedded"],setup_hint:"Set the root directory path and local port below.",service:be({name:"filesystem",type:"filesystem",path:"~/projects",port:4e3,description:"Embedded mcp-local filesystem server"})},{id:"filesystem-stdio",title:"Filesystem (Official @modelcontextprotocol/server-filesystem)",category:"Coding & Git",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Official Node.js stdio filesystem server with configurable allowed directories.",tags:["filesystem","files","stdio","official"],setup_hint:"Pass one or more allowed directory paths in Args.",service:be({name:"filesystem-stdio",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-filesystem","~/projects"],description:"Official stdio filesystem MCP server"})},{id:"playwright-mcp",title:"Playwright MCP (Official Microsoft Browser Automation)",category:"Browser & Automation",source:"github.com/mcp",source_url:"https://github.com/microsoft/playwright-mcp",popular:!0,description:"Official Microsoft Playwright MCP server using accessibility snapshots for fast, deterministic browser automation and testing.",tags:["playwright","microsoft","browser","automation","scraping","testing","headless"],service:be({name:"playwright",type:"stdio",command:"npx",args:["-y","@playwright/mcp@latest"],description:"Microsoft Playwright browser automation MCP"})},{id:"chrome-devtools-mcp",title:"Chrome DevTools MCP (Official Google DevTools)",category:"Browser & Automation",source:"github.com/mcp",source_url:"https://github.com/ChromeDevTools/chrome-devtools-mcp",popular:!0,description:"Inspect live DOM, network requests, console logs, performance traces, and debug web apps directly through Chrome DevTools.",tags:["chrome","devtools","debugging","network","performance","browser","google"],service:be({name:"chrome-devtools",type:"stdio",command:"npx",args:["-y","chrome-devtools-mcp@latest"],description:"Chrome DevTools inspection, network & performance trace MCP"})},{id:"puppeteer-mcp",title:"Puppeteer Browser Automation",category:"Browser & Automation",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Navigate pages, take screenshots, click elements, fill forms, and evaluate JavaScript in headless Chrome via Puppeteer.",tags:["puppeteer","browser","screenshots","chrome","automation"],service:be({name:"puppeteer",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-puppeteer"],description:"Puppeteer browser navigation & screenshot MCP"})},{id:"browserbase-mcp",title:"Browserbase & Stagehand Cloud Browser",category:"Browser & Automation",source:"github.com/mcp",source_url:"https://github.com/browserbase/mcp-server-browserbase",description:"Cloud browser automation, stealth web scraping, and Stagehand AI actions in managed headless browsers.",tags:["browserbase","stagehand","browser","cloud","scraping"],setup_hint:"Set BROWSERBASE_API_KEY and BROWSERBASE_PROJECT_ID in Environment Variables.",service:be({name:"browserbase",type:"stdio",command:"npx",args:["-y","@browserbasehq/mcp-server-browserbase"],env:{BROWSERBASE_API_KEY:"",BROWSERBASE_PROJECT_ID:""},description:"Browserbase cloud browser automation MCP"})},{id:"fetch-mcp",title:"Fetch — Web Page to Markdown Fetcher (Official)",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Official MCP server that fetches URLs and converts HTML pages into clean LLM-friendly Markdown.",tags:["fetch","web","markdown","scraper","http","url","official"],service:be({name:"fetch",type:"stdio",command:"uvx",args:["mcp-server-fetch"],description:"Official web content fetcher & HTML-to-Markdown converter"})},{id:"brave-search",title:"Brave Search MCP Server",category:"Search & Web",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Web search, news search, and local business search powered by the Brave Search API.",tags:["brave","search","web search","news","official"],setup_hint:"Requires BRAVE_API_KEY from https://brave.com/search/api/.",service:be({name:"brave-search",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-brave-search"],env:{BRAVE_API_KEY:""},description:"Brave Web & Local Search MCP"})},{id:"firecrawl-mcp",title:"Firecrawl — Web Scraping, Crawling & Search",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/mendableai/firecrawl-mcp-server",popular:!0,description:"Turn entire websites into LLM-ready markdown or structured data with deep crawling, batch scraping, and web search.",tags:["firecrawl","scraping","crawler","markdown","search","extract"],setup_hint:"Set FIRECRAWL_API_KEY in Environment Variables (or FIRECRAWL_API_URL for self-hosted).",service:be({name:"firecrawl",type:"stdio",command:"npx",args:["-y","firecrawl-mcp"],env:{FIRECRAWL_API_KEY:""},description:"Firecrawl web scraping, crawling & extraction MCP"})},{id:"tavily-mcp",title:"Tavily AI Search & Extract",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/tavily-ai/tavily-mcp",popular:!0,description:"Real-time web search and page extraction purpose-built for AI agents and RAG pipelines.",tags:["tavily","search","ai search","rag","extract","web"],setup_hint:"Requires TAVILY_API_KEY in Environment Variables.",service:be({name:"tavily",type:"stdio",command:"npx",args:["-y","tavily-mcp@latest"],env:{TAVILY_API_KEY:""},description:"Tavily AI real-time search & extraction MCP"})},{id:"exa-mcp",title:"Exa AI Semantic Search",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/exa-labs/exa-mcp-server",popular:!0,description:"Neural and keyword search across web pages, research papers, companies, and GitHub repositories.",tags:["exa","search","semantic search","research","papers"],setup_hint:"Requires EXA_API_KEY in Environment Variables.",service:be({name:"exa",type:"stdio",command:"npx",args:["-y","exa-mcp-server"],env:{EXA_API_KEY:""},description:"Exa neural web & research search MCP"})},{id:"perplexity-mcp",title:"Perplexity Ask / Sonar Research MCP",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/pplx-ai/modelcontextprotocol",popular:!0,description:"Live web research and cited answers powered by Perplexity Sonar API.",tags:["perplexity","sonar","search","research","citations"],setup_hint:"Requires PERPLEXITY_API_KEY in Environment Variables.",service:be({name:"perplexity",type:"stdio",command:"npx",args:["-y","@perplexity-ai/mcp-server"],env:{PERPLEXITY_API_KEY:""},description:"Perplexity Sonar web research MCP"})},{id:"brightdata-mcp",title:"Bright Data Web Unlocker & SERP MCP",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/luminati-io/brightdata-mcp",description:"Scrape protected websites, SERP results, and structured e-commerce data without blocks.",tags:["brightdata","scraping","proxy","serp","web"],setup_hint:"Requires API_TOKEN in Environment Variables.",service:be({name:"brightdata",type:"stdio",command:"npx",args:["-y","@brightdata/mcp"],env:{API_TOKEN:""},description:"Bright Data web scraping & SERP MCP"})},{id:"searxng-mcp",title:"SearXNG Self-Hosted Privacy Search",category:"Search & Web",source:"mcpservers.org",source_url:"https://github.com/ihor-sokoliuk/mcp-searxng",description:"Search the web privately through your own self-hosted SearXNG metasearch instance.",tags:["searxng","search","self-hosted","privacy","homelab"],setup_hint:"Set SEARXNG_URL to your SearXNG instance (e.g. http://localhost:8080).",service:be({name:"searxng",type:"stdio",command:"npx",args:["-y","mcp-searxng"],env:{SEARXNG_URL:"http://localhost:8080"},description:"Self-hosted SearXNG privacy web search MCP"})},{id:"apify-mcp",title:"Apify Actors Web Scraping & Automation",category:"Search & Web",source:"github.com/mcp",source_url:"https://github.com/apify/actors-mcp-server",description:"Run thousands of Apify Actors for social media, maps, e-commerce, and web data extraction.",tags:["apify","actors","scraping","crawler","extraction"],setup_hint:"Requires APIFY_TOKEN in Environment Variables.",service:be({name:"apify",type:"stdio",command:"npx",args:["-y","@apify/actors-mcp-server"],env:{APIFY_TOKEN:""},description:"Apify Actors scraping & automation MCP"})},{id:"youtube-transcript",title:"YouTube Transcript Extractor",category:"Search & Web",source:"mcpservers.org",source_url:"https://github.com/kimtaeyoon83/mcp-server-youtube-transcript",description:"Fetch subtitles and full transcripts from YouTube videos by URL or video ID.",tags:["youtube","transcript","captions","video","subtitles"],service:be({name:"youtube-transcript",type:"stdio",command:"npx",args:["-y","@kimtaeyoon83/mcp-server-youtube-transcript"],description:"YouTube video transcript & subtitle fetcher"})},{id:"postgres-mcp",title:"PostgreSQL MCP Server (Official)",category:"Databases & Storage",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Inspect database schemas, tables, and run read-only SQL queries against PostgreSQL.",tags:["postgres","postgresql","sql","database","official"],setup_hint:"Replace the connection URL in Args with your postgresql:// connection string.",service:be({name:"postgres",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-postgres","postgresql://localhost:5432/mydb"],description:"PostgreSQL schema inspection & SQL query MCP"})},{id:"dbhub-mcp",title:"DBHub — Universal Database Gateway (Postgres, MySQL, SQLite, SQL Server)",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/bytebase/dbhub",popular:!0,description:"Zero-dependency universal database MCP server supporting PostgreSQL, MySQL, MariaDB, SQLite, and SQL Server.",tags:["dbhub","bytebase","database","postgres","mysql","sqlite","mssql","sql"],setup_hint:"Pass --dsn with your database connection URL in Args.",service:be({name:"dbhub",type:"stdio",command:"npx",args:["-y","@bytebase/dbhub","--transport","stdio","--dsn","postgres://localhost:5432/mydb"],description:"Bytebase DBHub universal SQL database MCP"})},{id:"supabase-mcp",title:"Supabase MCP Server (Official)",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/supabase-community/supabase-mcp",popular:!0,description:"Manage Supabase projects, run SQL queries, inspect tables, manage migrations, and view edge function logs.",tags:["supabase","postgres","database","migrations","edge functions","baas"],setup_hint:"Set SUPABASE_ACCESS_TOKEN (Personal Access Token from Supabase dashboard).",service:be({name:"supabase",type:"stdio",command:"npx",args:["-y","@supabase/mcp-server-supabase@latest"],env:{SUPABASE_ACCESS_TOKEN:""},description:"Official Supabase database, auth & edge functions MCP"})},{id:"sqlite-mcp",title:"SQLite MCP Server",category:"Databases & Storage",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Query and analyze local SQLite database files, inspect schemas, and generate business insights.",tags:["sqlite","sql","database","local"],setup_hint:"Update --db-path in Args to point to your .db / .sqlite file.",service:be({name:"sqlite",type:"stdio",command:"uvx",args:["mcp-server-sqlite","--db-path","~/data/database.db"],description:"SQLite local database query & schema MCP"})},{id:"redis-mcp",title:"Redis MCP Server",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/redis/mcp-redis",description:"Interact with Redis keys, hashes, streams, pub/sub, and vector indexes.",tags:["redis","cache","kv","database","streams"],service:be({name:"redis",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-redis","redis://localhost:6379"],description:"Redis key-value & data structure MCP"})},{id:"neon-mcp",title:"Neon Serverless Postgres MCP",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/neondatabase/mcp-server-neon",description:"Create and manage Neon serverless Postgres branches, databases, and SQL migrations.",tags:["neon","postgres","serverless","database","branching"],setup_hint:"Requires NEON_API_KEY in Environment Variables.",service:be({name:"neon",type:"stdio",command:"npx",args:["-y","@neondatabase/mcp-server-neon","start"],env:{NEON_API_KEY:""},description:"Neon serverless Postgres management & SQL MCP"})},{id:"mongodb-mcp",title:"MongoDB Official MCP Server",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/mongodb-js/mongodb-mcp-server",description:"Query MongoDB collections, inspect document schemas, run aggregation pipelines, and manage Atlas clusters.",tags:["mongodb","nosql","atlas","database","collections"],setup_hint:"Set MDB_MCP_CONNECTION_STRING in Environment Variables.",service:be({name:"mongodb",type:"stdio",command:"npx",args:["-y","mongodb-mcp-server@latest"],env:{MDB_MCP_CONNECTION_STRING:"mongodb://localhost:27017/mydb"},description:"Official MongoDB & Atlas MCP server"})},{id:"qdrant-mcp",title:"Qdrant Vector Search & Memory MCP",category:"Databases & Storage",source:"github.com/mcp",source_url:"https://github.com/qdrant/mcp-server-qdrant",description:"Store and retrieve semantic memories and code snippets in a Qdrant vector database.",tags:["qdrant","vector","embeddings","rag","memory","database"],setup_hint:"Set QDRANT_URL and COLLECTION_NAME in Environment Variables.",service:be({name:"qdrant",type:"stdio",command:"uvx",args:["mcp-server-qdrant"],env:{QDRANT_URL:"http://localhost:6333",COLLECTION_NAME:"mcp-memory"},description:"Qdrant vector search & semantic memory MCP"})},{id:"notion-mcp",title:"Notion Official MCP Server",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://github.com/makenotion/notion-mcp-server",popular:!0,description:"Search Notion pages, query databases, create notes, and append blocks via the official Notion API.",tags:["notion","wiki","docs","notes","databases","official"],setup_hint:"Set OPENAPI_MCP_HEADERS with your Notion internal integration secret (ntn_...).",service:be({name:"notion",type:"stdio",command:"npx",args:["-y","@notionhq/notion-mcp-server"],env:{OPENAPI_MCP_HEADERS:'{"Authorization": "Bearer ntn_YOUR_NOTION_TOKEN", "Notion-Version": "2022-06-28"}'},description:"Official Notion pages & databases MCP server"})},{id:"linear-mcp",title:"Linear Official Remote MCP Server",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://linear.app/docs/mcp",popular:!0,description:"Search, create, and update Linear issues, projects, cycles, and comments via Linear hosted MCP.",tags:["linear","issues","project management","roadmap","remote","http"],setup_hint:"Add your Linear API key or OAuth Bearer token in the Authorization header.",service:be({name:"linear",type:"http",mcp_url:"https://mcp.linear.app/mcp",headers:{Authorization:"Bearer lin_api_YOUR_KEY"},description:"Official Linear issue tracking & project MCP"})},{id:"slack-mcp",title:"Slack MCP Server",category:"Productivity & Docs",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"List channels, post messages, reply to threads, add reactions, and read channel history in Slack workspaces.",tags:["slack","chat","messaging","channels","threads"],setup_hint:"Set SLACK_BOT_TOKEN (xoxb-...) and SLACK_TEAM_ID (T...) in Environment Variables.",service:be({name:"slack",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-slack"],env:{SLACK_BOT_TOKEN:"xoxb-",SLACK_TEAM_ID:""},description:"Slack workspace channels & messaging MCP"})},{id:"markitdown-mcp",title:"MarkItDown (Official Microsoft Document-to-Markdown)",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://github.com/microsoft/markitdown",popular:!0,description:"Convert PDF, PowerPoint, Word, Excel, images, audio, HTML, and ZIP files into clean Markdown.",tags:["markitdown","microsoft","pdf","docx","xlsx","pptx","markdown","converter"],service:be({name:"markitdown",type:"stdio",command:"uvx",args:["markitdown-mcp"],description:"Microsoft MarkItDown PDF/Office-to-Markdown converter MCP"})},{id:"figma-mcp",title:"Figma Context MCP (Framelink)",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://github.com/GLips/Figma-Context-MCP",popular:!0,description:"Give coding agents direct access to Figma layouts, frames, design tokens, and component hierarchies.",tags:["figma","design","ui","frontend","css","components"],setup_hint:"Set FIGMA_API_KEY in Environment Variables.",service:be({name:"figma",type:"stdio",command:"npx",args:["-y","figma-developer-mcp","--stdio"],env:{FIGMA_API_KEY:""},description:"Figma design layout & component context MCP"})},{id:"drawio-mcp",title:"draw.io Diagram Generator MCP",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://github.com/lgazo/drawio-mcp-server",description:"Create, edit, and inspect architecture diagrams and flowcharts in draw.io.",tags:["drawio","diagrams","architecture","flowchart"],service:be({name:"drawio",type:"stdio",command:"npx",args:["-y","@drawio/mcp"],description:"draw.io diagram generation & editing MCP"})},{id:"obsidian-mcp",title:"Obsidian Vault MCP Server",category:"Productivity & Docs",source:"mcpservers.org",source_url:"https://github.com/smithery-ai/mcp-obsidian",description:"Read, search, create, and link Markdown notes inside an Obsidian vault.",tags:["obsidian","notes","markdown","vault","pkm","knowledge base"],setup_hint:"Pass your Obsidian vault directory path in Args.",service:be({name:"obsidian",type:"stdio",command:"npx",args:["-y","mcp-obsidian","~/Documents/ObsidianVault"],description:"Obsidian Markdown vault search & note management MCP"})},{id:"atlassian-mcp",title:"Atlassian Remote MCP (Jira & Confluence)",category:"Productivity & Docs",source:"github.com/mcp",source_url:"https://www.atlassian.com/platform/remote-mcp-server",description:"Search and manage Jira issues, sprints, and Confluence pages through Atlassian Remote MCP.",tags:["atlassian","jira","confluence","issues","wiki","remote"],service:be({name:"atlassian",type:"http",mcp_url:"https://mcp.atlassian.com/v1/sse",headers:{Authorization:"Bearer YOUR_ATLASSIAN_TOKEN"},description:"Official Atlassian Jira & Confluence Remote MCP"})},{id:"google-maps-mcp",title:"Google Maps & Places MCP",category:"Productivity & Docs",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Geocoding, place search, directions, distance matrix, and elevation via Google Maps API.",tags:["google maps","maps","places","directions","geocoding"],setup_hint:"Requires GOOGLE_MAPS_API_KEY in Environment Variables.",service:be({name:"google-maps",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-google-maps"],env:{GOOGLE_MAPS_API_KEY:""},description:"Google Maps geocoding, directions & places MCP"})},{id:"time-mcp",title:"Time & Timezone Converter (Official)",category:"Productivity & Docs",source:"mcpservers.org",source_url:"https://github.com/modelcontextprotocol/servers",description:"Get current time and perform IANA timezone conversions.",tags:["time","timezone","clock","utc","official"],service:be({name:"time",type:"stdio",command:"uvx",args:["mcp-server-time"],description:"Official time & IANA timezone conversion MCP"})},{id:"sentry-mcp",title:"Sentry Error Tracking & Debugging MCP",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/getsentry/sentry-mcp",popular:!0,description:"Inspect Sentry issues, stacktraces, release health, and Seer root-cause analysis.",tags:["sentry","errors","debugging","stacktrace","observability","monitoring"],setup_hint:"Set SENTRY_ACCESS_TOKEN in Environment Variables.",service:be({name:"sentry",type:"stdio",command:"npx",args:["-y","@sentry/mcp-server@latest"],env:{SENTRY_ACCESS_TOKEN:""},description:"Sentry error tracking & stacktrace inspection MCP"})},{id:"stripe-mcp",title:"Stripe Official Agent Toolkit MCP",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/stripe/agent-toolkit",popular:!0,description:"Manage Stripe customers, products, prices, payment links, invoices, subscriptions, and refunds.",tags:["stripe","payments","billing","invoices","subscriptions","fintech"],setup_hint:"Set STRIPE_SECRET_KEY (rk_... or sk_test_...) in Environment Variables.",service:be({name:"stripe",type:"stdio",command:"npx",args:["-y","@stripe/mcp","--tools=all"],env:{STRIPE_SECRET_KEY:""},description:"Official Stripe payments, billing & subscriptions MCP"})},{id:"cloudflare-mcp",title:"Cloudflare Workers, KV, D1 & R2 MCP",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/cloudflare/mcp-server-cloudflare",popular:!0,description:"Manage Cloudflare Workers, KV namespaces, R2 buckets, D1 databases, and analytics.",tags:["cloudflare","workers","r2","d1","kv","cdn","dns"],service:be({name:"cloudflare",type:"stdio",command:"npx",args:["-y","@cloudflare/mcp-server-cloudflare"],description:"Cloudflare Workers, D1, R2 & KV management MCP"})},{id:"grafana-mcp",title:"Grafana Official Observability MCP",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/grafana/mcp-grafana",popular:!0,description:"Search dashboards, run PromQL and LogQL queries, inspect alert rules, and manage incidents in Grafana.",tags:["grafana","prometheus","loki","promql","logql","dashboards","observability","homelab"],setup_hint:"Set GRAFANA_URL and GRAFANA_API_KEY (service account token) in Environment Variables.",service:be({name:"grafana",type:"stdio",command:"uvx",args:["mcp-grafana"],env:{GRAFANA_URL:"http://localhost:3000",GRAFANA_API_KEY:""},description:"Official Grafana dashboards, PromQL, LogQL & alerts MCP"})},{id:"docker-mcp",title:"Docker Containers & Compose MCP",category:"Cloud & DevOps",source:"mcpservers.org",source_url:"https://github.com/ckreiling/mcp-server-docker",description:"List/start/stop Docker containers, inspect logs, manage images, volumes, and compose stacks.",tags:["docker","containers","compose","devops","logs"],service:be({name:"docker",type:"stdio",command:"uvx",args:["mcp-server-docker"],description:"Docker container, image & Compose management MCP"})},{id:"kubernetes-mcp",title:"Kubernetes Cluster (kubectl / Helm) MCP",category:"Cloud & DevOps",source:"mcpservers.org",source_url:"https://github.com/Flux159/mcp-server-kubernetes",description:"Inspect and manage Kubernetes pods, deployments, services, logs, events, and Helm releases.",tags:["kubernetes","k8s","kubectl","helm","pods","devops"],service:be({name:"kubernetes",type:"stdio",command:"npx",args:["-y","mcp-server-kubernetes"],description:"Kubernetes cluster & pod operations MCP"})},{id:"terraform-mcp",title:"HashiCorp Terraform Registry & IaC MCP",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/hashicorp/terraform-mcp-server",description:"Look up up-to-date Terraform Registry providers, resources, data sources, and module schemas.",tags:["terraform","hashicorp","iac","providers","modules","devops"],service:be({name:"terraform",type:"stdio",command:"docker",args:["run","-i","--rm","hashicorp/terraform-mcp-server"],description:"Official HashiCorp Terraform Registry & module docs MCP"})},{id:"azure-mcp",title:"Microsoft Azure Official MCP Server",category:"Cloud & DevOps",source:"github.com/mcp",source_url:"https://github.com/Azure/azure-mcp",description:"Manage Azure resources, Cosmos DB, Storage, Azure Monitor logs, Key Vault, and App Services.",tags:["azure","microsoft","cloud","cosmosdb","monitor"],service:be({name:"azure",type:"stdio",command:"npx",args:["-y","@azure/mcp@latest","server","start"],description:"Official Microsoft Azure cloud resources MCP"})},{id:"sequential-thinking",title:"Sequential Thinking — Structured Step-by-Step Reasoning",category:"AI & Memory",source:"github.com/mcp",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Official dynamic and reflective problem-solving tool through structured thought chains, branching, and revisions.",tags:["sequential-thinking","reasoning","planning","thinking","official"],service:be({name:"sequential-thinking",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-sequential-thinking"],description:"Structured step-by-step reflective reasoning MCP"})},{id:"memory-mcp",title:"Knowledge Graph Memory Server (Official)",category:"AI & Memory",source:"github.com/mcp",source_url:"https://github.com/modelcontextprotocol/servers",popular:!0,description:"Persistent memory using a local knowledge graph of entities, relations, and observations across sessions.",tags:["memory","knowledge graph","entities","persistent","official"],service:be({name:"memory",type:"stdio",command:"npx",args:["-y","@modelcontextprotocol/server-memory"],description:"Persistent local knowledge graph memory MCP"})},{id:"basic-memory",title:"Basic Memory — Local-First Markdown Knowledge Graph",category:"AI & Memory",source:"github.com/mcp",source_url:"https://github.com/basicmachines-co/basic-memory",description:"Persistent agent memory stored as clean local Markdown files compatible with Obsidian.",tags:["basic-memory","memory","markdown","obsidian","knowledge"],service:be({name:"basic-memory",type:"stdio",command:"uvx",args:["basic-memory","mcp"],description:"Local-first Markdown knowledge graph memory MCP"})},{id:"tpl-npx",title:"Blank npx (Node.js stdio) Template",category:"Templates",source:"template",description:"Starter template for any Node.js MCP package published to npm.",tags:["template","npx","node","npm","stdio"],service:be({name:"my-npx-mcp",type:"stdio",command:"npx",args:["-y","@scope/package-name"],description:"Custom npx stdio MCP server"})},{id:"tpl-uvx",title:"Blank uvx (Python stdio) Template",category:"Templates",source:"template",description:"Starter template for any Python MCP package published to PyPI and run via uvx.",tags:["template","uvx","python","pypi","stdio"],service:be({name:"my-uvx-mcp",type:"stdio",command:"uvx",args:["mcp-server-package"],description:"Custom uvx Python stdio MCP server"})},{id:"tpl-local-http",title:"Blank Local HTTP Daemon Template",category:"Templates",source:"template",description:"Starter template for a local background HTTP daemon supervised by mcp-local.",tags:["template","http","daemon","local"],service:be({name:"my-http-mcp",type:"http",command:"${HOME}/bin/my-mcp",port:7850,mcp_url:"http://localhost:7850/mcp",health_url:"http://localhost:7850/health",build_command:"make build",description:"Supervised local HTTP MCP daemon"})},{id:"tpl-remote-http",title:"Blank Remote Streamable HTTP SaaS Template",category:"Templates",source:"template",description:"Starter template for a hosted remote HTTP MCP endpoint with Bearer token auth.",tags:["template","remote","http","saas","bearer"],service:be({name:"remote-mcp",type:"http",mcp_url:"https://api.example.com/mcp",headers:{Authorization:"Bearer "},description:"Remote Streamable HTTP MCP server"})}];function xf(n){switch(n){case"zapier.com/mcp":return{background:"#ff4f00",color:"#ffffff"};case"github.com/mcp":return{background:"#000080",color:"#ffffff"};case"mcpservers.org":return{background:"#006644",color:"#ffffff"};case"template":return{background:"#555555",color:"#ffffff"};default:return{background:"#4b0082",color:"#ffffff"}}}function oy(n){switch(n){case"http":return{background:"#e0f0ff",color:"#003366",border:"1px solid #003366"};case"filesystem":return{background:"#fff4d6",color:"#664400",border:"1px solid #664400"};default:return{background:"#e8ffe8",color:"#005500",border:"1px solid #005500"}}}function iy(n){const l=n.service;return l.type==="filesystem"?`filesystem: ${l.path||"~/projects"} (port ${l.port||4e3})`:l.type==="http"&&!l.command?l.mcp_url||"https://...":[l.command,...l.args||[]].filter(Boolean).join(" ")||l.mcp_url||""}const ly=({selectedPreset:n,onSelect:l})=>{const[s,u]=m.useState(!1),[f,g]=m.useState(""),[h,v]=m.useState("all"),[x,E]=m.useState("All Categories"),[T,A]=m.useState(0),[P,U]=m.useState([]),[j,y]=m.useState(!1),$=m.useRef(null),Y=m.useRef(null),H=m.useRef(null);m.useEffect(()=>{const M=ne=>{$.current&&!$.current.contains(ne.target)&&u(!1)};return document.addEventListener("mousedown",M),()=>document.removeEventListener("mousedown",M)},[]),m.useEffect(()=>{const M=f.trim();if(M.length<2){U([]),y(!1);return}let ne=!1;y(!0);const de=window.setTimeout(async()=>{try{const Ae=await Le.searchRegistry(M);ne||U(Ae.presets||[])}catch{ne||U([])}finally{ne||y(!1)}},260);return()=>{ne=!0,window.clearTimeout(de)}},[f]);const V=m.useMemo(()=>{const M=f.trim().toLowerCase(),ne=M?M.split(/\s+/).filter(Boolean):[],de=W=>{if(h==="popular"&&!W.popular||h!=="all"&&h!=="popular"&&W.source!==h||x!=="All Categories"&&W.category!==x)return!1;if(ne.length===0)return!0;const ce=[W.title,W.service.name,W.category,W.source,W.description,W.service.command,...W.service.args||[],W.service.mcp_url,...W.tags||[]].filter(Boolean).join(" ").toLowerCase();return ne.every(xe=>ce.includes(xe))},Ae=ms.filter(de),$e=new Set(Ae.map(W=>W.service.name.toLowerCase())),De=[];if(ne.length>0&&(h==="all"||h==="github.com/mcp")&&x==="All Categories")for(const W of P){const ce=W.service.name.toLowerCase();$e.has(ce)||($e.add(ce),De.push({...W,source:"MCP Registry Live"}))}return[...Ae,...De]},[f,h,x,P]);m.useEffect(()=>{A(0)},[f,h,x]);const J=M=>{const ne={...M,service:{...M.service,args:[...M.service.args||[]],env:{...M.service.env||{}},headers:{...M.service.headers||{}},deps:[...M.service.deps||[]],extra_ports:[...M.service.extra_ports||[]]}};l(ne),u(!1)},L=M=>{if(!s&&(M.key==="ArrowDown"||M.key==="Enter")){u(!0);return}s&&(M.key==="ArrowDown"?(M.preventDefault(),A(ne=>Math.min(ne+1,Math.max(0,V.length-1)))):M.key==="ArrowUp"?(M.preventDefault(),A(ne=>Math.max(0,ne-1))):M.key==="Enter"?V[T]&&(M.preventDefault(),J(V[T])):M.key==="Escape"&&(M.preventDefault(),u(!1)))};return c.jsx(Ct,{label:"Popular Service Presets (github.com/mcp • mcpservers.org • zapier.com/mcp)",children:c.jsxs("div",{ref:$,style:{position:"relative",display:"flex",flexDirection:"column",gap:8},children:[c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center",flexWrap:"wrap"},children:[c.jsxs(he,{type:"button",onClick:()=>{const M=!s;u(M),M&&window.setTimeout(()=>{var ne;return(ne=Y.current)==null?void 0:ne.focus()},20)},active:s,style:{fontWeight:"bold",minWidth:235,justifyContent:"space-between"},children:[c.jsxs("span",{children:["📋 ",n?`Preset: ${n.service.name}`:"Browse Popular Presets"]}),c.jsx("span",{style:{marginLeft:8},children:s?"▴":"▾"})]}),c.jsxs("div",{style:{flex:1,minWidth:240,display:"flex",gap:4,alignItems:"center"},children:[c.jsx(Ke,{ref:Y,value:f,placeholder:"Find as you type: zapier, playwright, github, postgres, context7, grafana, notion...",onFocus:()=>u(!0),onChange:M=>{g(M.target.value),s||u(!0)},onKeyDown:L,fullWidth:!0}),f&&c.jsx(he,{size:"sm",type:"button",onClick:()=>{var M;g(""),(M=Y.current)==null||M.focus()},title:"Clear search",children:"✕"})]})]}),c.jsxs("div",{style:{display:"flex",gap:4,flexWrap:"wrap",alignItems:"center"},children:[c.jsx("span",{style:{fontSize:11,color:"#333",marginRight:2},children:"Quick pick:"}),["zapier-mcp-all","github-stdio","playwright-mcp","context7-stdio","fetch-mcp","postgres-mcp","brave-search","sequential-thinking","filesystem-embedded"].map(M=>{const ne=ms.find(Ae=>Ae.id===M);if(!ne)return null;const de=(n==null?void 0:n.id)===ne.id;return c.jsx(he,{size:"sm",type:"button",active:de,onClick:()=>J(ne),style:{fontSize:11,height:24,padding:"0 6px"},title:ne.description,children:ne.service.name},ne.id)}),c.jsxs(he,{size:"sm",type:"button",onClick:()=>{u(!0),window.setTimeout(()=>{var M;return(M=Y.current)==null?void 0:M.focus()},20)},style:{fontSize:11,height:24,padding:"0 8px",fontWeight:"bold"},children:["All (",ms.length,"+) ▾"]})]}),s&&c.jsxs("div",{style:{position:"absolute",top:40,left:0,right:0,zIndex:100,background:"#c0c0c0",borderTop:"2px solid #ffffff",borderLeft:"2px solid #ffffff",borderRight:"2px solid #000000",borderBottom:"2px solid #000000",boxShadow:"4px 4px 0px rgba(0,0,0,0.45)",padding:8,display:"flex",flexDirection:"column",gap:6},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",gap:6,flexWrap:"wrap"},children:[c.jsx("div",{style:{display:"flex",gap:4,flexWrap:"wrap"},children:ny.map(M=>c.jsx(he,{size:"sm",type:"button",active:h===M.id,onClick:()=>v(M.id),style:{fontSize:11,height:24,padding:"0 7px"},children:M.label},M.id))}),c.jsx("div",{style:{fontSize:11,color:"#222"},children:j?"Searching live MCP Registry...":`${V.length} matching service${V.length===1?"":"s"}`})]}),c.jsx("div",{style:{display:"flex",gap:4,flexWrap:"wrap",paddingBottom:2,borderBottom:"1px solid #808080"},children:ry.map(M=>c.jsx("button",{type:"button",onClick:()=>E(M),style:{fontSize:11,padding:"2px 7px",cursor:"pointer",border:x===M?"1px solid #000080":"1px solid #808080",background:x===M?"#000080":"#dfdfdf",color:x===M?"#ffffff":"#000000",fontWeight:x===M?"bold":"normal"},children:M},M))}),c.jsx("div",{ref:H,style:{maxHeight:310,overflowY:"auto",background:"#ffffff",borderTop:"2px solid #808080",borderLeft:"2px solid #808080",borderRight:"2px solid #dfdfdf",borderBottom:"2px solid #dfdfdf"},children:V.length===0?c.jsx("div",{style:{padding:16,textAlign:"center",fontSize:12,color:"#444"},children:j?"Searching MCP Registry...":`No presets matched "${f}". Try clearing filters or searching another keyword.`}):V.map((M,ne)=>{const de=ne===T,Ae=iy(M);return c.jsxs("div",{onMouseEnter:()=>A(ne),onClick:()=>J(M),style:{padding:"7px 10px",cursor:"pointer",borderBottom:"1px solid #e4e4e4",background:de?"#000080":ne%2===0?"#ffffff":"#f7f7f7",color:de?"#ffffff":"#000000",display:"flex",flexDirection:"column",gap:3},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",gap:8},children:[c.jsxs("div",{style:{display:"flex",alignItems:"center",gap:6,flexWrap:"wrap"},children:[c.jsxs("span",{style:{fontWeight:"bold",fontSize:13},children:[M.popular?"★ ":"",M.title]}),c.jsx("span",{style:{fontSize:10,padding:"1px 5px",fontFamily:"monospace",...de?{background:"#ffffff",color:"#000080"}:oy(M.service.type||"stdio")},children:(M.service.type||"stdio").toUpperCase()}),c.jsx("span",{style:{fontSize:10,padding:"1px 5px",...xf(M.source)},children:M.source})]}),c.jsx("span",{style:{fontSize:11,opacity:.85,whiteSpace:"nowrap"},children:M.category})]}),c.jsx("div",{style:{fontSize:11,color:de?"#e8e8ff":"#333333",lineHeight:1.3},children:M.description}),Ae&&c.jsxs("div",{style:{fontSize:11,fontFamily:"monospace",color:de?"#ffff99":"#004080",overflow:"hidden",textOverflow:"ellipsis",whiteSpace:"nowrap"},children:["▸ ",Ae]})]},M.id)})}),c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",fontSize:11},children:[c.jsxs("span",{style:{color:"#333"},children:["Tip: Use ",c.jsx("b",{children:"↑/↓"})," + ",c.jsx("b",{children:"Enter"})," to select, or type any package name to search the live MCP Registry."]}),c.jsx(he,{size:"sm",type:"button",onClick:()=>u(!1),children:"Close ▴"})]})]}),n&&c.jsxs("div",{style:{background:"#ffffe1",border:"1px solid #808000",padding:"6px 10px",fontSize:12,display:"flex",flexDirection:"column",gap:3},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",gap:8},children:[c.jsxs("div",{children:[c.jsx("b",{children:"Loaded Preset:"})," ",n.title," ",c.jsx("span",{style:{fontSize:10,padding:"1px 5px",marginLeft:4,...xf(n.source)},children:n.source})]}),n.source_url&&c.jsx("a",{href:n.source_url,target:"_blank",rel:"noreferrer",style:{color:"#000080",fontWeight:"bold",fontSize:11},children:"Open Docs / Repo ↗"})]}),n.setup_hint&&c.jsxs("div",{style:{color:"#443000",fontSize:11},children:["💡 ",c.jsx("b",{children:"Setup Note:"})," ",n.setup_hint]})]})]})})},ay=({onCreate:n,onCancel:l})=>{const[s,u]=m.useState(()=>ep()),[f,g]=m.useState(null),[h,v]=m.useState(!0),[x,E]=m.useState(!1);return c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsx(ly,{selectedPreset:f,onSelect:T=>{g(T),u(T.service)}}),c.jsx(ut,{checked:h,onChange:T=>v(T.target.checked),label:"Start & register with enabled agents immediately after creating"}),c.jsx(tp,{initial:s,isNew:!0,saving:x,onCancel:l,onSubmit:async T=>{E(!0);try{await n(T,h)}finally{E(!1)}}})]})},sy=({agents:n,onNotify:l,onAfterMutate:s})=>{const[u,f]=m.useState(n),[g,h]=m.useState([]),[v,x]=m.useState(0),[E,T]=m.useState(!1);m.useEffect(()=>{f(n)},[n]);const A=m.useCallback(async()=>{try{const y=await Le.agentConfigs();h(y.configs||[])}catch(y){l(y instanceof Error?y.message:String(y),!0)}},[l]);m.useEffect(()=>{A()},[A]);const P=async()=>{T(!0);try{await Le.updateAgents(u),l("Saved agent registration targets to config.yaml."),await s(),await A()}catch(y){l(y instanceof Error?y.message:String(y),!0)}finally{T(!1)}},U=async()=>{T(!0);try{const y=await Le.registerAll();l(y.output||(y.ok?"Registered all services.":y.error||"Failed"),!y.ok),await s(),await A()}catch(y){l(y instanceof Error?y.message:String(y),!0)}finally{T(!1)}},j=g[v];return c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsx(Ct,{label:"Managed Agent Targets (config.yaml -> agents:)",children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:8},children:[c.jsx(ut,{checked:u.managed,onChange:y=>f({...u,managed:y.target.checked}),label:"Explicitly manage agents: key (when unchecked, legacy mode enables all three)"}),c.jsxs("div",{style:{display:"flex",gap:18,flexWrap:"wrap",paddingLeft:8},children:[c.jsx(ut,{checked:u.managed?u.opencode:!0,disabled:!u.managed,onChange:y=>f({...u,opencode:y.target.checked}),label:"OpenCode"}),c.jsx(ut,{checked:u.managed?u.cursor:!0,disabled:!u.managed,onChange:y=>f({...u,cursor:y.target.checked}),label:"Cursor"}),c.jsx(ut,{checked:u.managed?u.claude:!0,disabled:!u.managed,onChange:y=>f({...u,claude:y.target.checked}),label:"Claude Desktop (uses mcp-local bridge for HTTP)"})]}),c.jsxs("div",{style:{display:"flex",gap:8,justifyContent:"flex-end",flexWrap:"wrap"},children:[c.jsx(he,{size:"sm",onClick:A,disabled:E,children:"Refresh Files"}),c.jsx(he,{size:"sm",onClick:U,disabled:E,children:"Register All Enabled Services"}),c.jsx(he,{size:"sm",primary:!0,onClick:P,disabled:E,children:"Save Agent Targets"})]})]})}),g.length>0&&c.jsxs("div",{children:[c.jsx(ll,{value:v,onChange:y=>x(y),children:g.map((y,$)=>c.jsxs(ar,{value:$,children:[y.host," ",y.exists?"":"(missing)"]},y.host))}),c.jsx(il,{style:{padding:10},children:j&&c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:6},children:[c.jsxs("div",{style:{fontSize:12,fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"Path:"})," ",j.path," | ",c.jsx("b",{children:"Enabled:"})," ",j.enabled?"yes":"no"]}),j.error?c.jsx("div",{style:{color:"#c00000",fontSize:12},children:j.error}):c.jsx(wr,{style:{height:240,background:"#fff",color:"#000",padding:8},children:c.jsx("pre",{style:{margin:0,fontFamily:"monospace",fontSize:12},children:j.exists?j.content:"(File does not exist yet on disk)"})})]})})]})]})},gs=({size:n=32})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"2",y:"1",width:"12",height:"9",fill:"#c0c0c0",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"3",y:"2",width:"10",height:"7",fill:"#008080"}),c.jsx("rect",{x:"4",y:"3",width:"3",height:"1",fill:"#00ffff"}),c.jsx("rect",{x:"4",y:"5",width:"6",height:"1",fill:"#ffff00"}),c.jsx("rect",{x:"6",y:"10",width:"4",height:"2",fill:"#808080"}),c.jsx("rect",{x:"1",y:"12",width:"14",height:"3",fill:"#c0c0c0",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"10",y:"13",width:"3",height:"1",fill:"#008000"})]}),vs=({size:n=32})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"1",y:"2",width:"6",height:"5",fill:"#c0c0c0",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"2",y:"3",width:"4",height:"3",fill:"#000080"}),c.jsx("rect",{x:"9",y:"2",width:"6",height:"5",fill:"#c0c0c0",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"10",y:"3",width:"4",height:"3",fill:"#008080"}),c.jsx("rect",{x:"3",y:"7",width:"2",height:"3",fill:"#808080"}),c.jsx("rect",{x:"11",y:"7",width:"2",height:"3",fill:"#808080"}),c.jsx("rect",{x:"2",y:"10",width:"12",height:"2",fill:"#008000"}),c.jsx("rect",{x:"7",y:"12",width:"2",height:"2",fill:"#008000"}),c.jsx("rect",{x:"4",y:"14",width:"8",height:"1",fill:"#000"})]}),ys=({size:n=32})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"2",y:"4",width:"12",height:"10",fill:"#ffffff",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"5",y:"2",width:"6",height:"2",fill:"#808080",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"7",y:"6",width:"2",height:"6",fill:"#ff0000"}),c.jsx("rect",{x:"5",y:"8",width:"6",height:"2",fill:"#ff0000"})]}),xs=({size:n=32})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"1",y:"1",width:"14",height:"14",fill:"#000080",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"3",y:"1",width:"9",height:"5",fill:"#c0c0c0"}),c.jsx("rect",{x:"8",y:"2",width:"3",height:"3",fill:"#000080"}),c.jsx("rect",{x:"3",y:"8",width:"10",height:"6",fill:"#ffffff"}),c.jsx("rect",{x:"4",y:"10",width:"8",height:"1",fill:"#808080"}),c.jsx("rect",{x:"4",y:"12",width:"6",height:"1",fill:"#808080"})]}),ws=({size:n=32})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"2",y:"1",width:"12",height:"14",fill:"#ffffff",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"2",y:"1",width:"12",height:"3",fill:"#008080"}),c.jsx("rect",{x:"4",y:"6",width:"8",height:"1",fill:"#000"}),c.jsx("rect",{x:"4",y:"8",width:"6",height:"1",fill:"#000080"}),c.jsx("rect",{x:"4",y:"10",width:"8",height:"1",fill:"#000"}),c.jsx("rect",{x:"4",y:"12",width:"5",height:"1",fill:"#808080"})]}),wf=({size:n=16})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"6",y:"1",width:"4",height:"14",fill:"#808080"}),c.jsx("rect",{x:"1",y:"6",width:"14",height:"4",fill:"#808080"}),c.jsx("rect",{x:"3",y:"3",width:"10",height:"10",fill:"#c0c0c0",stroke:"#000",strokeWidth:"1"}),c.jsx("rect",{x:"6",y:"6",width:"4",height:"4",fill:"#008080",stroke:"#000",strokeWidth:"1"})]}),cy=({size:n=16})=>c.jsxs("svg",{width:n,height:n,viewBox:"0 0 16 16",shapeRendering:"crispEdges",children:[c.jsx("rect",{x:"2",y:"2",width:"5",height:"5",fill:"#ff0000"}),c.jsx("rect",{x:"9",y:"2",width:"5",height:"5",fill:"#00a800"}),c.jsx("rect",{x:"2",y:"9",width:"5",height:"5",fill:"#0000ff"}),c.jsx("rect",{x:"9",y:"9",width:"5",height:"5",fill:"#ffff00"})]}),Ps=({state:n})=>{let l="#808080",s=n;switch(n){case"running":l="#00b000";break;case"degraded":case"orphaned":l="#d8a000";break;case"port conflict":case"stale":l="#e00000";break;case"stopped":l="#808080";break;default:l="#808080"}return c.jsxs("span",{style:{display:"inline-flex",alignItems:"center",gap:6,whiteSpace:"nowrap"},children:[c.jsx("span",{"aria-hidden":"true",style:{width:9,height:9,display:"inline-block",background:l,border:"1px solid #000",boxShadow:"inset 1px 1px 0 rgba(255,255,255,0.55)",flexShrink:0}}),c.jsx("span",{children:s})]})},uy=({onNotify:n,onAfterMutate:l})=>{const[s,u]=m.useState(null),[f,g]=m.useState(null),[h,v]=m.useState(!1),x=m.useCallback(async(E=!1)=>{v(!0);try{const[T,A]=await Promise.all([Le.validate(),Le.doctor("",E)]);u(T),g(A),E&&(n(A.remaining===0?"Doctor --fix resolved all fixable issues.":`Doctor --fix finished with ${A.remaining} remaining issue(s).`,A.remaining>0),await l())}catch(T){n(T instanceof Error?T.message:String(T),!0)}finally{v(!1)}},[n,l]);return m.useEffect(()=>{x(!1)},[x]),c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsxs("div",{style:{display:"flex",gap:8,flexWrap:"wrap",justifyContent:"space-between",alignItems:"center"},children:[c.jsx("div",{style:{fontSize:12},children:"Checks config paths, binaries, dependencies, stale PIDs, orphan processes, port conflicts, and agent registrations."}),c.jsxs("div",{style:{display:"flex",gap:6},children:[c.jsx(he,{size:"sm",onClick:()=>x(!1),disabled:h,children:"Run Checks"}),c.jsx(he,{size:"sm",primary:!0,onClick:()=>x(!0),disabled:h,children:h?"Running...":"Doctor --fix"})]})]}),s&&c.jsx(Ct,{label:`Config Validation (${s.valid?"Valid":"Errors Found"})`,children:s.issues.length===0?c.jsx("div",{style:{fontSize:12,color:"#006400",fontWeight:"bold"},children:"No configuration errors or warnings."}):c.jsx("div",{style:{display:"flex",flexDirection:"column",gap:4},children:s.issues.map((E,T)=>c.jsxs(an,{variant:"well",style:{padding:"4px 8px",fontSize:12,background:E.level==="error"?"#ffe0e0":"#ffffd0",color:"#000"},children:[c.jsxs("b",{children:["[",E.level.toUpperCase(),"] ",E.service,":"]})," ",E.message]},T))})}),f&&c.jsxs(Ct,{label:`Doctor Findings (${f.remaining} remaining)`,children:[c.jsx("div",{style:{maxHeight:220,overflow:"auto",border:"2px inset #808080",background:"#fff"},children:c.jsxs(Lo,{children:[c.jsx(Io,{children:c.jsxs(Ht,{children:[c.jsx(tt,{style:{width:150},children:"Service"}),c.jsx(tt,{style:{width:120},children:"State"}),c.jsx(tt,{children:"Findings"})]})}),c.jsx(zo,{children:f.services.map(E=>c.jsxs(Ht,{children:[c.jsx(et,{style:{fontWeight:"bold"},children:E.name}),c.jsx(et,{children:c.jsx(Ps,{state:E.state})}),c.jsx(et,{children:E.findings.length===0?c.jsx("span",{style:{color:"#006400"},children:"OK — no problems"}):c.jsx("ul",{style:{margin:0,paddingLeft:16},children:E.findings.map((T,A)=>c.jsxs("li",{children:[T.fixed?"✅ Fixed: ":T.fixable?"🔧 ":"⚠️ ",T.message,T.error?` (error: ${T.error})`:""]},A))})})]},E.name))})]})}),f.output&&c.jsx(wr,{style:{maxHeight:130,marginTop:8,background:"#000",color:"#c0ffc0",padding:8},children:c.jsx("pre",{style:{margin:0,fontFamily:"monospace",fontSize:11,whiteSpace:"pre-wrap"},children:f.output})})]})]})},dy=({onNotify:n,onAfterMutate:l})=>{const[s,u]=m.useState([]),[f,g]=m.useState({}),[h,v]=m.useState(!1),[x,E]=m.useState(!1),T=m.useCallback(async()=>{E(!0);try{const U=(await Le.importPreview()).candidates||[];u(U);const j={};for(const y of U)y.exists||(j[`${y.source}:${y.name}`]=!0);g(j)}catch(P){n(P instanceof Error?P.message:String(P),!0)}finally{E(!1)}},[n]);m.useEffect(()=>{T()},[T]);const A=async()=>{const P=s.filter(U=>f[`${U.source}:${U.name}`]).map(U=>({source:U.source,name:U.name}));if(P.length===0){n("Select at least one service to import.",!0);return}E(!0);try{const U=await Le.importApply(P,h);n(`Imported ${U.imported} service(s) into config.yaml.`),await l(),await T()}catch(U){n(U instanceof Error?U.message:String(U),!0)}finally{E(!1)}};return c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",flexWrap:"wrap",gap:8},children:[c.jsx("span",{style:{fontSize:12},children:"Discover MCP server entries already defined in OpenCode, Cursor, and Claude Desktop configs."}),c.jsx(he,{size:"sm",onClick:T,disabled:x,children:"Rescan Agents"})]}),c.jsx("div",{style:{maxHeight:280,overflow:"auto",border:"2px inset #808080",background:"#fff"},children:c.jsxs(Lo,{children:[c.jsx(Io,{children:c.jsxs(Ht,{children:[c.jsx(tt,{style:{width:50},children:"Pick"}),c.jsx(tt,{style:{width:90},children:"Source"}),c.jsx(tt,{children:"Name"}),c.jsx(tt,{style:{width:80},children:"Mode"}),c.jsx(tt,{children:"Command / Endpoint"}),c.jsx(tt,{style:{width:100},children:"Status"})]})}),c.jsx(zo,{children:s.length===0?c.jsx(Ht,{children:c.jsx("td",{colSpan:6,style:{textAlign:"center",padding:16},children:x?"Scanning agent configs...":"No MCP entries found in OpenCode, Cursor, or Claude Desktop."})}):s.map(P=>{const U=`${P.source}:${P.name}`;return c.jsxs(Ht,{children:[c.jsx(et,{style:{textAlign:"center"},children:c.jsx(ut,{checked:!!f[U],onChange:j=>g({...f,[U]:j.target.checked})})}),c.jsx(et,{children:P.source}),c.jsx(et,{style:{fontWeight:"bold"},children:P.name}),c.jsx(et,{children:P.mode}),c.jsx(et,{style:{fontFamily:"monospace",fontSize:12},children:P.mcp_url||[P.command,...P.args||[]].filter(Boolean).join(" ")}),c.jsx(et,{children:P.exists?"Already in config":"New"})]},U)})})]})}),c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",flexWrap:"wrap",gap:8},children:[c.jsx(ut,{checked:h,onChange:P=>v(P.target.checked),label:"Overwrite services that already exist in config.yaml"}),c.jsx(he,{primary:!0,onClick:A,disabled:x||s.length===0,children:"Import Selected"})]})]})},fy=({service:n,onNotify:l})=>{const[s,u]=m.useState([]),[f,g]=m.useState(n.resolved_log_path),[h,v]=m.useState(!0),[x,E]=m.useState(!1),[T,A]=m.useState(!0),[P,U]=m.useState(!1),j=m.useRef(null),y=m.useCallback(async()=>{U(!0);try{const $=await Le.getLogs(n.name,250);g($.path),v($.exists),u($.lines||[])}catch($){l($ instanceof Error?$.message:String($),!0)}finally{U(!1)}},[n.name,l]);return m.useEffect(()=>{x||y()},[x,y]),m.useEffect(()=>{if(!x)return;u([]);const $=new EventSource(`/api/services/${encodeURIComponent(n.name)}/logs/stream`);return $.addEventListener("log",Y=>{v(!0),u(H=>{const V=[...H,Y.data];return V.length>1e3?V.slice(V.length-1e3):V})}),$.addEventListener("info",Y=>{u(H=>[...H,`[info] ${Y.data}`])}),$.onerror=()=>{},()=>$.close()},[x,n.name]),m.useEffect(()=>{if(T&&j.current){const $=j.current.parentElement;$&&($.scrollTop=$.scrollHeight)}},[s,T]),c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:8},children:[c.jsxs("div",{style:{display:"flex",gap:8,flexWrap:"wrap",alignItems:"center",justifyContent:"space-between"},children:[c.jsxs("span",{style:{fontSize:12,fontFamily:"monospace",wordBreak:"break-all"},children:[f," ",!h&&"(file not created yet)"]}),c.jsxs("div",{style:{display:"flex",gap:8,alignItems:"center",flexWrap:"wrap"},children:[c.jsx(ut,{checked:x,onChange:$=>E($.target.checked),label:"Follow Live (SSE)"}),c.jsx(ut,{checked:T,onChange:$=>A($.target.checked),label:"Auto-scroll"}),c.jsx(he,{size:"sm",onClick:y,disabled:P||x,children:"Refresh"}),c.jsx(he,{size:"sm",onClick:()=>u([]),children:"Clear"}),c.jsx(he,{size:"sm",onClick:()=>{var $;($=navigator.clipboard)==null||$.writeText(s.join(`
`)),l("Copied log lines to clipboard.")},children:"Copy"})]})]}),c.jsx(wr,{style:{height:290,background:"#000000",color:"#c0ffc0",padding:8},children:c.jsx("pre",{ref:j,style:{margin:0,fontFamily:"ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",fontSize:12,lineHeight:1.4,whiteSpace:"pre-wrap",wordBreak:"break-word"},children:s.length>0?s.join(`
`):P?"Reading log file...":x?"Waiting for log output...":"(empty log)"})})]})},py=({onNotify:n,onAfterMutate:l})=>{const[s,u]=m.useState(null),[f,g]=m.useState(""),[h,v]=m.useState(!1),[x,E]=m.useState(!1),T=m.useCallback(async()=>{E(!0);try{const P=await Le.getRawConfig();u(P),g(P.content||""),v(!1)}catch(P){n(P instanceof Error?P.message:String(P),!0)}finally{E(!1)}},[n]);m.useEffect(()=>{T()},[T]);const A=async()=>{E(!0);try{const P=await Le.putRawConfig(f);u(P),g(P.content),v(!1),n("Saved config.yaml (with automatic backup)."),await l()}catch(P){n(P instanceof Error?P.message:String(P),!0)}finally{E(!1)}};return c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsxs("div",{style:{display:"flex",justifyContent:"space-between",alignItems:"center",flexWrap:"wrap",gap:8},children:[c.jsxs("span",{style:{fontSize:12,fontFamily:"monospace",wordBreak:"break-all"},children:[(s==null?void 0:s.path)||"~/.mcp-local/config.yaml"," ",h?" (modified)":""]}),c.jsxs("div",{style:{display:"flex",gap:6},children:[c.jsx(he,{size:"sm",onClick:T,disabled:x,children:"Reload"}),c.jsx(he,{size:"sm",primary:!0,onClick:A,disabled:x,children:"Validate & Save YAML"})]})]}),c.jsx(Ke,{multiline:!0,rows:16,value:f,onChange:P=>{g(P.target.value),v(!0)},placeholder:"# Paste or edit ~/.mcp-local/config.yaml here",fullWidth:!0,style:{fontFamily:"ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",fontSize:12}}),s&&s.backups&&s.backups.length>0&&c.jsx(Ct,{label:`Recent Backups (${s.backups.length})`,children:c.jsx(wr,{style:{maxHeight:90,background:"#fff",color:"#000",padding:6},children:s.backups.map(P=>c.jsxs("div",{style:{fontSize:11,fontFamily:"monospace",marginBottom:2},children:[P.timestamp," (",P.size," bytes) — ",P.path]},P.path))})})]})},hy=[{value:"",label:"(default / complete)"},{value:"core",label:"core"},{value:"extended",label:"extended"},{value:"complete",label:"complete"}],my=[{value:"",label:"(default)"},{value:"core",label:"core"},{value:"extended",label:"extended"},{value:"complete",label:"complete"}],gy=({service:n,onNotify:l,onAfterMutate:s})=>{const[u,f]=m.useState(null),[g,h]=m.useState(""),[v,x]=m.useState(!1),[E,T]=m.useState(!1),[A,P]=m.useState(""),[U,j]=m.useState([]),[y,$]=m.useState(""),[Y,H]=m.useState(!1),[V,J]=m.useState(!1),[L,M]=m.useState(!1),[ne,de]=m.useState(!1),[Ae,$e]=m.useState(!1),De=m.useCallback(async()=>{H(!0);try{const K=await Le.getTools(n.name);f(K),h(K.active_tier||""),x(K.code_mode),T(K.no_code_mode),P(n.tools_config_path||""),j(K.tools||[]),de(!1)}catch(K){l(K instanceof Error?K.message:String(K),!0)}finally{H(!1)}},[n.name,n.tools_config_path,l]);m.useEffect(()=>{De()},[De]);const W=m.useMemo(()=>{const K=y.trim().toLowerCase();return K?U.map((ue,ge)=>({t:ue,idx:ge})).filter(({t:ue})=>ue.name.toLowerCase().includes(K)||(ue.description||"").toLowerCase().includes(K)||(ue.effective_tier||"").toLowerCase().includes(K)):U.map((ue,ge)=>({t:ue,idx:ge}))},[U,y]),ce=async()=>{J(!0);try{const K=await Le.syncTools(n.name);f(K),j(K.tools||[]),de(!1),l(`Synced ${K.tools.length} tool(s) from live server for ${n.name}.`),await s()}catch(K){l(K instanceof Error?K.message:String(K),!0)}finally{J(!1)}},xe=async K=>{J(!0);try{const ue=await Le.updateTools(n.name,{active_tier:g,code_mode:v,no_code_mode:E,tools_config_path:A,tools:U,apply:!0});if(f(ue),j(ue.tools||[]),de(!1),K){const ge=await Le.restartService(n.name);$e(!1),l(ge.ok?`Saved tools & restarted ${n.name}.`:ge.error||"Restart failed",!ge.ok)}else n.running&&$e(!0),l(`Saved ${n.name} tool policy (${ue.overrides_count} override(s) written to ${ue.tools_config_path}).`);await s()}catch(ue){l(ue instanceof Error?ue.message:String(ue),!0)}finally{J(!1)}};return c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsx(Ct,{label:"Global Tier & Code Mode Policy",children:c.jsxs("div",{style:{display:"flex",gap:12,flexWrap:"wrap",alignItems:"flex-end"},children:[c.jsxs("div",{style:{minWidth:180},children:[c.jsx("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:"Active Tier (AST_MCP_TIER)"}),c.jsx(Nn,{options:hy,value:g,onChange:K=>{h(K.value),de(!0)},width:180})]}),c.jsxs("div",{style:{minWidth:220,flex:1},children:[c.jsxs("label",{style:{display:"block",fontSize:12,marginBottom:4,fontWeight:"bold"},children:["tools.json Path (resolved: ",(u==null?void 0:u.tools_config_path)||n.resolved_tools_path,")"]}),c.jsx(Ke,{value:A,onChange:K=>{P(K.target.value),de(!0)},placeholder:"~/.astcache/tools.json",fullWidth:!0})]}),c.jsxs("div",{style:{display:"flex",gap:12,alignItems:"center",paddingBottom:4},children:[c.jsx(ut,{checked:v,onChange:K=>{x(K.target.checked),K.target.checked&&T(!1),de(!0)},label:"Code Mode"}),c.jsx(ut,{checked:E,onChange:K=>{T(K.target.checked),K.target.checked&&x(!1),de(!0)},label:"No Code Mode"})]})]})}),Ae&&c.jsxs(an,{variant:"well",style:{padding:"6px 10px",background:"#ffffd0",color:"#000",fontSize:12},children:["Note: ",c.jsx("b",{children:n.name})," loads ",c.jsx("code",{children:"tools.json"})," and ",c.jsx("code",{children:"AST_MCP_*"})," env vars at startup only. Restart to apply changes live."]}),c.jsxs("div",{style:{display:"flex",gap:8,flexWrap:"wrap",alignItems:"center",justifyContent:"space-between"},children:[c.jsxs("div",{style:{display:"flex",gap:6,flexWrap:"wrap"},children:[c.jsx(he,{size:"sm",onClick:ce,disabled:V||!n.mcp_endpoint,children:"Sync from Live Server"}),c.jsx(he,{size:"sm",onClick:()=>M(!L),children:L?"Hide tools.json":"Preview tools.json"}),c.jsx(he,{size:"sm",onClick:De,disabled:Y,children:"Reload"})]}),c.jsxs("div",{style:{display:"flex",gap:6,alignItems:"center"},children:[c.jsx(Ke,{value:y,onChange:K=>$(K.target.value),placeholder:"Filter tools...",style:{width:180}}),c.jsx(he,{size:"sm",primary:!0,onClick:()=>xe(!1),disabled:V,children:ne?"Save & Write tools.json *":"Save & Write tools.json"}),n.running&&c.jsx(he,{size:"sm",onClick:()=>xe(!0),disabled:V,children:"Save & Restart"})]})]}),L&&u&&c.jsx(Ct,{label:`tools.json Overrides (${u.overrides_count})`,children:c.jsx(wr,{style:{maxHeight:160,background:"#fff",color:"#000",padding:8},children:c.jsx("pre",{style:{margin:0,fontFamily:"monospace",fontSize:12},children:u.overrides_json||"{}"})})}),c.jsx("div",{style:{maxHeight:300,overflow:"auto",border:"2px inset #808080",background:"#fff"},children:c.jsxs(Lo,{children:[c.jsx(Io,{children:c.jsxs(Ht,{children:[c.jsx(tt,{style:{width:64},children:"Enabled"}),c.jsx(tt,{children:"Tool Name"}),c.jsx(tt,{style:{width:120},children:"Tier Override"}),c.jsx(tt,{style:{width:95},children:"Effective"}),c.jsx(tt,{children:"Description Override"})]})}),c.jsx(zo,{children:W.length===0?c.jsx(Ht,{children:c.jsx("td",{colSpan:5,style:{textAlign:"center",padding:16},children:Y?"Loading tools...":U.length===0?'No tools configured in config.yaml. Click "Sync from Live Server" while the service is running.':"No tools match your filter."})}):W.map(({t:K,idx:ue})=>c.jsxs(Ht,{children:[c.jsx(et,{style:{textAlign:"center"},children:c.jsx(ut,{checked:K.enabled,onChange:ge=>{const F=[...U];F[ue]={...K,enabled:ge.target.checked},j(F),de(!0)}})}),c.jsx(et,{style:{fontFamily:"monospace",fontWeight:"bold"},children:K.name}),c.jsx(et,{children:c.jsx(Gf,{options:my,value:K.tier||"",onChange:ge=>{const F=[...U];F[ue]={...K,tier:ge.value},j(F),de(!0)},width:"100%"})}),c.jsx(et,{style:{fontSize:12},children:K.tier||K.effective_tier||"—"}),c.jsx(et,{children:c.jsx(Ke,{value:K.description||"",onChange:ge=>{const F=[...U];F[ue]={...K,description:ge.target.value},j(F),de(!0)},placeholder:"(server default)",fullWidth:!0})})]},K.name))})]})})]})},on=({title:n,icon:l,bounds:s,zIndex:u,active:f,minimized:g,maximized:h,onFocus:v,onClose:x,onMinimize:E,onToggleMaximize:T,onMove:A,children:P,footer:U})=>{const j=m.useRef(null),[y,$]=m.useState(!1);if(g)return null;const Y=V=>{if(V.target.closest("button")||(v(),h))return;j.current={startX:V.clientX,startY:V.clientY,origX:s.x,origY:s.y},$(!0);const J=M=>{if(!j.current)return;const ne=Math.max(0,j.current.origX+(M.clientX-j.current.startX)),de=Math.max(0,j.current.origY+(M.clientY-j.current.startY));A(ne,de)},L=()=>{j.current=null,$(!1),window.removeEventListener("mousemove",J),window.removeEventListener("mouseup",L)};window.addEventListener("mousemove",J),window.addEventListener("mouseup",L)},H=h?{position:"fixed",top:0,left:0,right:0,bottom:48,width:"100%",height:"calc(100vh - 48px)",zIndex:u,display:"flex",flexDirection:"column"}:{position:"fixed",top:s.y,left:s.x,width:`min(${s.width}px, calc(100vw - 16px))`,height:s.height?`min(${s.height}px, calc(100vh - 64px))`:void 0,maxHeight:"calc(100vh - 60px)",zIndex:u,display:"flex",flexDirection:"column"};return c.jsxs(Ws,{style:H,onMouseDown:v,children:[c.jsxs(Fs,{active:f,onMouseDown:Y,style:{display:"flex",alignItems:"center",justifyContent:"space-between",cursor:h?"default":y?"grabbing":"grab",userSelect:"none",gap:8,flexShrink:0},children:[c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:6,overflow:"hidden",textOverflow:"ellipsis",whiteSpace:"nowrap"},children:[l,c.jsx("span",{style:{overflow:"hidden",textOverflow:"ellipsis"},children:n})]}),c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:2,flexShrink:0},children:[E&&c.jsx(he,{size:"sm",square:!0,onClick:E,title:"Minimize",style:{minWidth:22,padding:0,fontWeight:"bold"},children:"_"}),T&&c.jsx(he,{size:"sm",square:!0,onClick:T,title:h?"Restore":"Maximize",style:{minWidth:22,padding:0,fontWeight:"bold"},children:h?"❐":"□"}),c.jsx(he,{size:"sm",square:!0,onClick:x,title:"Close",style:{minWidth:22,padding:0,fontWeight:"bold"},children:"×"})]})]}),c.jsx(Bs,{style:{flex:1,minHeight:0,display:"flex",flexDirection:"column",overflow:"auto",padding:10},children:P}),U]})},Xi={original:{label:"Windows 95 Classic",theme:_0,bg:"#008080"},rainyDay:{label:"Rainy Day",theme:L0,bg:"#3e5368"},vaporTeal:{label:"Vapor Teal",theme:D0,bg:"#016868"},lilac:{label:"Lilac",theme:V0,bg:"#5a4e7c"},rose:{label:"Rose",theme:K0,bg:"#734a54"},tokyoDark:{label:"Tokyo Dark",theme:P0,bg:"#16161e"},modernDark:{label:"Modern Dark",theme:N0,bg:"#1e1e1e"},matrix:{label:"Matrix",theme:W0,bg:"#081808"},hotdogStand:{label:"Hotdog Stand",theme:X0,bg:"#b00000"}},vy=Xm`
  ${am}
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${S0}') format('woff2');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${C0}') format('woff2');
    font-weight: bold;
    font-style: normal;
  }
  html, body, #root {
    height: 100%;
    margin: 0;
    padding: 0;
    overflow: hidden;
    font-family: 'ms_sans_serif', sans-serif;
    background: ${n=>n.$bg};
  }
  * {
    box-sizing: border-box;
  }
  code, pre {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  }
`;function bf(n){return n.toLocaleTimeString([],{hour:"2-digit",minute:"2-digit"})}function yy(){var ot,Fe;const[n,l]=m.useState(()=>{const N=localStorage.getItem("mcp-local-95-theme");return N&&N in Xi?N:"original"}),[s,u]=m.useState(null),[f,g]=m.useState({config_exists:!1,agents:{managed:!1,opencode:!0,cursor:!0,claude:!0},services:[]}),[h,v]=m.useState(""),[x,E]=m.useState(0),[T,A]=m.useState(""),[P,U]=m.useState(!0),[j,y]=m.useState(!1),[$,Y]=m.useState(""),[H,V]=m.useState(null),[J,L]=m.useState(null),[M,ne]=m.useState(!1),[de,Ae]=m.useState(!1),[$e,De]=m.useState(()=>bf(new Date)),[W,ce]=m.useState(""),[xe,K]=m.useState(!0),[ue,ge]=m.useState(20),[F,se]=m.useState("manager"),[G,S]=m.useState({manager:{open:!0,minimized:!1,maximized:!1,z:10,bounds:{x:108,y:14,width:1040,height:700}},add:{open:!1,minimized:!1,maximized:!1,z:11,bounds:{x:150,y:30,width:820}},agents:{open:!1,minimized:!1,maximized:!1,z:12,bounds:{x:190,y:55,width:780}},doctor:{open:!1,minimized:!1,maximized:!1,z:13,bounds:{x:210,y:70,width:780}},import:{open:!1,minimized:!1,maximized:!1,z:14,bounds:{x:230,y:85,width:760}},raw:{open:!1,minimized:!1,maximized:!1,z:15,bounds:{x:200,y:45,width:820}},delete:{open:!1,minimized:!1,maximized:!1,z:19,bounds:{x:340,y:180,width:420}}}),C=m.useCallback((N,te=!1)=>{L({msg:N,isError:te})},[]);m.useEffect(()=>{if(!J)return;const N=setTimeout(()=>L(null),6e3);return()=>clearTimeout(N)},[J]),m.useEffect(()=>{const N=setInterval(()=>De(bf(new Date)),15e3);return()=>clearInterval(N)},[]);const I=m.useCallback(async()=>{try{const[N,te]=await Promise.all([Le.info(),Le.state()]);u(N),g(te),v(Ee=>{var ve;return Ee&&te.services.some(Ue=>Ue.name===Ee)?Ee:((ve=te.services[0])==null?void 0:ve.name)||""})}catch(N){C(N instanceof Error?N.message:String(N),!0)}},[C]);m.useEffect(()=>{I()},[I]),m.useEffect(()=>{if(!P)return;const N=setInterval(()=>{Le.state().then(g).catch(()=>{})},5e3);return()=>clearInterval(N)},[P]);const D=m.useCallback(N=>{const te=ue+1;ge(te),se(N),S(Ee=>({...Ee,[N]:{...Ee[N],open:!0,minimized:!1,z:te}}))},[ue]),q=N=>{S(te=>({...te,[N]:{...te[N],open:!1,minimized:!1}}))},b=N=>{S(te=>({...te,[N]:{...te[N],minimized:!0}}))},oe=N=>{S(te=>({...te,[N]:{...te[N],maximized:!te[N].maximized}})),D(N)},Q=(N,te,Ee)=>{S(ve=>({...ve,[N]:{...ve[N],bounds:{...ve[N].bounds,x:te,y:Ee}}}))},_=m.useMemo(()=>f.services.find(N=>N.name===h),[f.services,h]),le=m.useMemo(()=>_?Z0(_):null,[_]);m.useEffect(()=>{V(null)},[h]);const fe=m.useMemo(()=>{const N=T.trim().toLowerCase();return N?f.services.filter(te=>te.name.toLowerCase().includes(N)||te.mode.toLowerCase().includes(N)||te.state.toLowerCase().includes(N)||(te.description||"").toLowerCase().includes(N)||(te.command||"").toLowerCase().includes(N)||(te.mcp_url||"").toLowerCase().includes(N)):f.services},[f.services,T]),ae=async(N,te)=>{y(!0);try{const Ee=await te();Y(Ee.output||Ee.error||N),C(Ee.ok?`${N}: OK`:Ee.error||`${N} failed`,!Ee.ok),await I()}catch(Ee){const ve=Ee instanceof Error?Ee.message:String(Ee);Y(ve),C(ve,!0)}finally{y(!1)}},je=async N=>{y(!0);try{const te=await Le.healthService(N);V(te),C(te.ok?`Health OK (${te.status_code}, ${te.latency_ms}ms)`:te.error||`Health HTTP ${te.status_code}`,!te.ok)}catch(te){C(te instanceof Error?te.message:String(te),!0)}finally{y(!1)}},ze=f.services.filter(N=>N.running).length,Ce=Xi[n],Me={manager:"MCP Control Panel",add:"Add New MCP Service",agents:"Agent Sync & IDE Configs",doctor:"Doctor & Config Validator",import:"Import MCP Servers",raw:"Notepad — config.yaml",delete:"Confirm Delete"},Re=[{id:"manager",label:"MCP Manager",icon:c.jsx(gs,{size:32})},{id:"agents",label:"Agent Sync",icon:c.jsx(vs,{size:32})},{id:"doctor",label:"Doctor",icon:c.jsx(ys,{size:32})},{id:"import",label:"Import IDEs",icon:c.jsx(xs,{size:32})},{id:"raw",label:"config.yaml",icon:c.jsx(ws,{size:32})}],Ne=(N,te,Ee)=>{var ur;if(!(!f.agents.managed||te==="opencode"&&f.agents.opencode||te==="cursor"&&f.agents.cursor||te==="claude"&&f.agents.claude))return c.jsx("span",{title:`${te}: disabled in agents config`,style:{opacity:.35,fontSize:11},children:Ee},te);const Ue=!!((ur=N.registered)!=null&&ur[te]);return c.jsx("span",{title:`${te}: ${Ue?"registered":"not registered"}`,style:{fontSize:10,fontWeight:"bold",padding:"1px 4px",border:"1px solid #000",background:Ue?"#c8f7c5":"#ffd6d6",color:"#000"},children:Ee},te)};return c.jsxs(Ym,{theme:Ce.theme,children:[c.jsx(vy,{$bg:Ce.bg}),c.jsxs("div",{style:{position:"fixed",inset:0,bottom:48,overflow:"hidden"},onMouseDown:()=>{M&&(ne(!1),Ae(!1))},children:[c.jsx("div",{style:{position:"absolute",top:14,left:12,display:"flex",flexDirection:"column",gap:16,zIndex:1},children:Re.map(N=>c.jsxs("button",{type:"button",onClick:()=>D(N.id),style:{width:78,background:"transparent",border:"none",color:"#ffffff",textShadow:"1px 1px 0 #000000",display:"flex",flexDirection:"column",alignItems:"center",gap:4,cursor:"pointer",padding:4,fontSize:12,fontFamily:"inherit"},children:[N.icon,c.jsx("span",{children:N.label})]},N.id))}),G.manager.open&&c.jsx(on,{id:"manager",title:`MCP Control Panel — ${(s==null?void 0:s.config_path)||"~/.mcp-local/config.yaml"}`,icon:c.jsx(gs,{size:16}),bounds:G.manager.bounds,zIndex:G.manager.z,active:F==="manager",minimized:G.manager.minimized,maximized:G.manager.maximized,onFocus:()=>D("manager"),onClose:()=>q("manager"),onMinimize:()=>b("manager"),onToggleMaximize:()=>oe("manager"),onMove:(N,te)=>Q("manager",N,te),footer:c.jsxs(an,{variant:"well",style:{padding:"3px 8px",fontSize:12,display:"flex",justifyContent:"space-between",gap:12,flexShrink:0},children:[c.jsx("span",{children:J?c.jsx("span",{style:{color:J.isError?"#c00000":"#006400",fontWeight:"bold"},children:J.msg}):c.jsxs("span",{children:[f.services.length," service(s) | ",ze," running",s!=null&&s.config_is_symlink?` | symlink -> ${s.config_symlink_target}`:""]})}),c.jsxs("span",{children:["v",(s==null?void 0:s.version)||"dev"]})]}),children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10,height:"100%"},children:[c.jsxs($o,{style:{flexWrap:"wrap",gap:6,justifyContent:"space-between"},children:[c.jsxs("div",{style:{display:"flex",gap:4,flexWrap:"wrap",alignItems:"center"},children:[c.jsx(he,{size:"sm",primary:!0,onClick:()=>D("add"),children:"+ New Service"}),c.jsx(he,{size:"sm",onClick:()=>ae("Start All",Le.startAll),disabled:j,children:"Start All"}),c.jsx(he,{size:"sm",onClick:()=>ae("Stop All",()=>Le.stopAll(!0)),disabled:j,children:"Stop All"}),c.jsx(he,{size:"sm",onClick:()=>ae("Restart All",Le.restartAll),disabled:j,children:"Restart All"}),c.jsx(he,{size:"sm",onClick:()=>ae("Register All",Le.registerAll),disabled:j,children:"Register All"}),c.jsx(he,{size:"sm",onClick:()=>D("agents"),children:"Agents"}),c.jsx(he,{size:"sm",onClick:()=>D("doctor"),children:"Doctor"}),c.jsx(he,{size:"sm",onClick:()=>D("import"),children:"Import"}),c.jsx(he,{size:"sm",onClick:()=>D("raw"),children:"YAML"})]}),c.jsxs("div",{style:{display:"flex",gap:8,alignItems:"center"},children:[c.jsx(Ke,{value:T,onChange:N=>A(N.target.value),placeholder:"Search services...",style:{width:150}}),c.jsx(ut,{checked:P,onChange:N=>U(N.target.checked),label:"Poll 5s"}),c.jsx(he,{size:"sm",onClick:I,disabled:j,children:"Refresh"})]})]}),f.config_error&&c.jsxs(an,{variant:"well",style:{padding:8,background:"#ffe0e0",color:"#900000",fontSize:12},children:[c.jsx("b",{children:"Config Error:"})," ",f.config_error," ",c.jsx(he,{size:"sm",onClick:()=>D("raw"),style:{marginLeft:8},children:"Open config.yaml"})]}),c.jsx("div",{style:{minHeight:150,maxHeight:225,overflow:"auto",border:"2px inset #808080",background:"#fff",flexShrink:0},children:c.jsxs(Lo,{children:[c.jsx(Io,{children:c.jsxs(Ht,{children:[c.jsx(tt,{children:"Service"}),c.jsx(tt,{style:{width:90},children:"Mode"}),c.jsx(tt,{style:{width:125},children:"Status"}),c.jsx(tt,{style:{width:70},children:"Port"}),c.jsx(tt,{style:{width:75},children:"PID"}),c.jsx(tt,{style:{width:118},children:"Agents"}),c.jsx(tt,{style:{width:65},children:"Tools"}),c.jsx(tt,{style:{width:160},children:"Quick Actions"})]})}),c.jsx(zo,{children:fe.length===0?c.jsx(Ht,{children:c.jsx("td",{colSpan:8,style:{textAlign:"center",padding:20},children:f.services.length===0?c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:8,alignItems:"center"},children:[c.jsxs("span",{children:["No MCP services configured yet",f.config_exists?".":" (~/.mcp-local/config.yaml will be created on save)."]}),c.jsxs("div",{style:{display:"flex",gap:8},children:[c.jsx(he,{size:"sm",primary:!0,onClick:()=>D("add"),children:"+ Add First Service"}),c.jsx(he,{size:"sm",onClick:()=>D("import"),children:"Import from IDEs"})]})]}):"No services match your search."})}):fe.map(N=>{var Ee;const te=N.name===h;return c.jsxs(Ht,{onClick:()=>v(N.name),style:{cursor:"pointer",background:te?"#000080":void 0,color:te?"#ffffff":void 0},children:[c.jsx(et,{style:{fontWeight:"bold"},children:N.name}),c.jsx(et,{children:N.mode}),c.jsx(et,{children:c.jsx(Ps,{state:N.state})}),c.jsx(et,{children:N.port||"—"}),c.jsx(et,{children:N.pid||"—"}),c.jsx(et,{children:c.jsxs("div",{style:{display:"flex",gap:3},children:[Ne(N,"opencode","OC"),Ne(N,"cursor","CR"),Ne(N,"claude","CL")]})}),c.jsx(et,{children:((Ee=N.tools)==null?void 0:Ee.length)||0}),c.jsx(et,{onClick:ve=>ve.stopPropagation(),children:c.jsxs("div",{style:{display:"flex",gap:4},children:[N.running?c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Stop ${N.name}`,()=>Le.stopService(N.name,!0)),children:"Stop"}):c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Start ${N.name}`,()=>Le.startService(N.name,!0)),children:"Start"}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Restart ${N.name}`,()=>Le.restartService(N.name)),children:"Restart"})]})})]},N.name)})})]})}),_?c.jsxs("div",{style:{flex:1,minHeight:0,display:"flex",flexDirection:"column"},children:[c.jsxs(ll,{value:x,onChange:N=>E(N),children:[c.jsx(ar,{value:0,children:"Overview & Controls"}),c.jsx(ar,{value:1,children:"Edit Configuration"}),c.jsxs(ar,{value:2,children:["Tools & Tiers (",((ot=_.tools)==null?void 0:ot.length)||0,")"]}),c.jsx(ar,{value:3,children:"Logs"})]}),c.jsxs(il,{style:{flex:1,overflow:"auto",padding:10},children:[x===0&&c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:10},children:[c.jsxs("div",{style:{display:"flex",gap:6,flexWrap:"wrap",alignItems:"center",justifyContent:"space-between"},children:[c.jsxs("div",{style:{display:"flex",gap:6,flexWrap:"wrap"},children:[c.jsx(he,{size:"sm",primary:!0,disabled:j,onClick:()=>ae(`Start ${_.name}`,()=>Le.startService(_.name,!0)),children:"Start"}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Stop ${_.name}`,()=>Le.stopService(_.name,!0)),children:"Stop"}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Restart ${_.name}`,()=>Le.restartService(_.name)),children:"Restart"}),_.effective_build&&c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Rebuild ${_.name}`,()=>Le.rebuildService(_.name)),children:"Rebuild"}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Register ${_.name}`,()=>Le.registerService(_.name)),children:"Register Agents"}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>ae(`Deregister ${_.name}`,()=>Le.deregisterService(_.name)),children:"Deregister"}),_.health_url&&c.jsx(he,{size:"sm",disabled:j,onClick:()=>je(_.name),children:"Test Health"}),_.dashboard_url&&c.jsx(he,{size:"sm",onClick:()=>window.open(_.dashboard_url,"_blank","noopener,noreferrer"),children:"Open Dashboard ↗"})]}),c.jsx(he,{size:"sm",disabled:j,onClick:()=>{ce(_.name),K(!0),D("delete")},children:"Delete Service"})]}),c.jsx(Ct,{label:`Service Details: ${_.name}`,children:c.jsxs("div",{style:{display:"grid",gridTemplateColumns:"repeat(auto-fit, minmax(260px, 1fr))",gap:"6px 16px",fontSize:12},children:[c.jsxs("div",{children:[c.jsx("b",{children:"Mode:"})," ",_.mode," | ",c.jsx("b",{children:"Status:"})," ",c.jsx(Ps,{state:_.state})]}),c.jsxs("div",{children:[c.jsx("b",{children:"Detail:"})," ",_.detail]}),c.jsxs("div",{children:[c.jsx("b",{children:"PID / Uptime:"})," ",_.pid?`${_.pid}${_.uptime?` (${_.uptime})`:""}`:"—"]}),c.jsxs("div",{children:[c.jsx("b",{children:"Ports:"})," ",_.port||"—",(Fe=_.extra_ports)!=null&&Fe.length?` (extra: ${_.extra_ports.join(", ")})`:""]}),c.jsxs("div",{style:{gridColumn:"1 / -1",fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"Command:"})," ",_.path?`embedded filesystem (${_.path})`:_.command?[_.command,..._.args||[]].join(" "):"(remote-only)"]}),c.jsxs("div",{style:{fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"MCP Endpoint:"})," ",_.mcp_endpoint||"—"]}),c.jsxs("div",{style:{fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"Health URL:"})," ",_.health_url||"—"]}),c.jsxs("div",{style:{fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"Log Path:"})," ",_.resolved_log_path]}),c.jsxs("div",{style:{fontFamily:"monospace",wordBreak:"break-all"},children:[c.jsx("b",{children:"Installer:"})," ",_.effective_installer||"native (mcp-local)"]}),c.jsxs("div",{children:[c.jsx("b",{children:"Active Tier:"})," ",_.active_tier||"(default)"," |"," ",c.jsx("b",{children:"Tool Overrides:"})," ",_.overrides_count]}),c.jsxs("div",{children:[c.jsx("b",{children:"Registered:"})," ",Object.entries(_.registered||{}).map(([N,te])=>`${N}: ${te?"yes":"no"}`).join(" | ")]})]})}),H&&c.jsxs(Ct,{label:`Health Probe (${H.ok?"Healthy":"Failed"} — ${H.latency_ms}ms)`,children:[c.jsxs("div",{style:{fontSize:12,marginBottom:4},children:[c.jsx("b",{children:"URL:"})," ",c.jsx("code",{children:H.url||"—"})," ",H.status_code?`| HTTP ${H.status_code}`:"",H.error?` | Error: ${H.error}`:""]}),H.body&&c.jsx(wr,{style:{maxHeight:110,background:"#fff",color:"#000",padding:6},children:c.jsx("pre",{style:{margin:0,fontSize:11},children:H.body})})]}),$&&c.jsx(Ct,{label:"Last Action Output",children:c.jsx(wr,{style:{maxHeight:110,background:"#000",color:"#c0ffc0",padding:6},children:c.jsx("pre",{style:{margin:0,fontSize:11,whiteSpace:"pre-wrap"},children:$})})})]}),x===1&&le&&c.jsx(tp,{initial:le,saving:j,onSubmit:async N=>{y(!0);try{const te=await Le.updateService(_.name,N);v(te.name),C(`Saved changes to ${te.name}.`),await I()}catch(te){C(te instanceof Error?te.message:String(te),!0)}finally{y(!1)}}}),x===2&&c.jsx(gy,{service:_,onNotify:C,onAfterMutate:I}),x===3&&c.jsx(fy,{service:_,onNotify:C})]})]}):c.jsx(an,{variant:"well",style:{padding:16,textAlign:"center",fontSize:12},children:"Select a service above to inspect controls, edit configuration, manage tool tiers, or view logs."})]})}),G.add.open&&c.jsx(on,{id:"add",title:"Add New MCP Service Wizard",icon:c.jsx(wf,{size:16}),bounds:G.add.bounds,zIndex:G.add.z,active:F==="add",minimized:G.add.minimized,maximized:G.add.maximized,onFocus:()=>D("add"),onClose:()=>q("add"),onMinimize:()=>b("add"),onToggleMaximize:()=>oe("add"),onMove:(N,te)=>Q("add",N,te),children:c.jsx(ay,{onCancel:()=>q("add"),onCreate:async(N,te)=>{try{const Ee=await Le.createService(N);if(v(Ee.name),q("add"),te){const ve=await Le.startService(Ee.name,!0);Y(ve.output),C(ve.ok?`Created, started, and registered ${Ee.name}.`:`Created ${Ee.name}, start reported: ${ve.error}`,!ve.ok)}else C(`Created service ${Ee.name}.`);await I()}catch(Ee){C(Ee instanceof Error?Ee.message:String(Ee),!0)}}})}),G.agents.open&&c.jsx(on,{id:"agents",title:"Agent Sync & IDE Configs",icon:c.jsx(vs,{size:16}),bounds:G.agents.bounds,zIndex:G.agents.z,active:F==="agents",minimized:G.agents.minimized,maximized:G.agents.maximized,onFocus:()=>D("agents"),onClose:()=>q("agents"),onMinimize:()=>b("agents"),onToggleMaximize:()=>oe("agents"),onMove:(N,te)=>Q("agents",N,te),children:c.jsx(sy,{agents:f.agents,onNotify:C,onAfterMutate:I})}),G.doctor.open&&c.jsx(on,{id:"doctor",title:"MCP Doctor & Config Validator",icon:c.jsx(ys,{size:16}),bounds:G.doctor.bounds,zIndex:G.doctor.z,active:F==="doctor",minimized:G.doctor.minimized,maximized:G.doctor.maximized,onFocus:()=>D("doctor"),onClose:()=>q("doctor"),onMinimize:()=>b("doctor"),onToggleMaximize:()=>oe("doctor"),onMove:(N,te)=>Q("doctor",N,te),children:c.jsx(uy,{onNotify:C,onAfterMutate:I})}),G.import.open&&c.jsx(on,{id:"import",title:"Import MCP Servers from Agent IDEs",icon:c.jsx(xs,{size:16}),bounds:G.import.bounds,zIndex:G.import.z,active:F==="import",minimized:G.import.minimized,maximized:G.import.maximized,onFocus:()=>D("import"),onClose:()=>q("import"),onMinimize:()=>b("import"),onToggleMaximize:()=>oe("import"),onMove:(N,te)=>Q("import",N,te),children:c.jsx(dy,{onNotify:C,onAfterMutate:I})}),G.raw.open&&c.jsx(on,{id:"raw",title:"Notepad — config.yaml",icon:c.jsx(ws,{size:16}),bounds:G.raw.bounds,zIndex:G.raw.z,active:F==="raw",minimized:G.raw.minimized,maximized:G.raw.maximized,onFocus:()=>D("raw"),onClose:()=>q("raw"),onMinimize:()=>b("raw"),onToggleMaximize:()=>oe("raw"),onMove:(N,te)=>Q("raw",N,te),children:c.jsx(py,{onNotify:C,onAfterMutate:I})}),G.delete.open&&W&&c.jsx(on,{id:"delete",title:`Confirm Delete — ${W}`,bounds:G.delete.bounds,zIndex:G.delete.z,active:F==="delete",onFocus:()=>D("delete"),onClose:()=>q("delete"),onMove:(N,te)=>Q("delete",N,te),children:c.jsxs("div",{style:{display:"flex",flexDirection:"column",gap:12},children:[c.jsxs("div",{style:{fontSize:13},children:["Remove ",c.jsx("b",{children:W})," from ",c.jsx("code",{children:"config.yaml"}),"?"]}),c.jsx(ut,{checked:xe,onChange:N=>K(N.target.checked),label:"Also deregister from OpenCode, Cursor, and Claude Desktop"}),c.jsxs("div",{style:{display:"flex",justifyContent:"flex-end",gap:8},children:[c.jsx(he,{onClick:()=>q("delete"),children:"Cancel"}),c.jsx(he,{primary:!0,disabled:j,onClick:async()=>{const N=W;q("delete"),await ae(`Delete ${N}`,()=>Le.deleteService(N,xe))},children:"Delete"})]})]})})]}),M&&c.jsxs("div",{style:{position:"fixed",bottom:46,left:4,zIndex:9999,display:"flex",alignItems:"flex-end"},children:[c.jsxs(Es,{style:{minWidth:215},children:[c.jsx(lr,{onClick:()=>{D("manager"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(gs,{size:16})," MCP Control Panel"]})}),c.jsx(lr,{onClick:()=>{D("add"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(wf,{size:16})," New MCP Service..."]})}),c.jsx(Po,{}),c.jsx(lr,{onClick:()=>{D("agents"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(vs,{size:16})," Agent Sync & IDEs"]})}),c.jsx(lr,{onClick:()=>{D("doctor"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(ys,{size:16})," Doctor & Validator"]})}),c.jsx(lr,{onClick:()=>{D("import"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(xs,{size:16})," Import from Agents..."]})}),c.jsx(lr,{onClick:()=>{D("raw"),ne(!1)},children:c.jsxs("span",{style:{display:"flex",alignItems:"center",gap:8},children:[c.jsx(ws,{size:16})," Edit config.yaml"]})}),c.jsx(Po,{}),c.jsx(lr,{onClick:()=>Ae(!de),children:c.jsxs("span",{children:["Theme: ",Ce.label," ▸"]})}),c.jsx(lr,{onClick:()=>{ne(!1),I()},children:c.jsx("span",{children:"Refresh All"})})]}),de&&c.jsx(Es,{style:{marginLeft:2,minWidth:185},children:Object.keys(Xi).map(N=>c.jsx(lr,{onClick:()=>{l(N),localStorage.setItem("mcp-local-95-theme",N),Ae(!1),ne(!1)},children:c.jsxs("span",{children:[N===n?"✓ ":"   ",Xi[N].label]})},N))})]}),c.jsx(If,{style:{top:"auto",bottom:0,height:48,zIndex:9998},children:c.jsxs($o,{style:{justifyContent:"space-between",height:"100%",padding:"0 6px"},children:[c.jsxs("div",{style:{display:"flex",alignItems:"center",gap:6,flex:1,minWidth:0,overflow:"hidden"},children:[c.jsxs(he,{onClick:N=>{N.stopPropagation(),ne(!M),Ae(!1)},active:M,style:{fontWeight:"bold",display:"flex",alignItems:"center",gap:6,flexShrink:0},children:[c.jsx(cy,{size:16}),"Start"]}),c.jsx(Po,{orientation:"vertical",size:28}),c.jsx("div",{style:{display:"flex",gap:4,overflow:"hidden"},children:Object.keys(G).filter(N=>G[N].open).map(N=>{const te=F===N&&!G[N].minimized;return c.jsx(he,{active:te,onClick:()=>{te?b(N):D(N)},style:{maxWidth:175,overflow:"hidden",textOverflow:"ellipsis",whiteSpace:"nowrap",fontWeight:te?"bold":"normal"},children:Me[N]},N)})})]}),c.jsxs(an,{variant:"well",style:{padding:"4px 10px",display:"flex",alignItems:"center",gap:10,fontSize:12,flexShrink:0},children:[c.jsxs("span",{title:"Running services / total configured",children:[ze,"/",f.services.length," UP"]}),c.jsx("span",{children:$e})]})]})})]})}lm.createRoot(document.getElementById("root")).render(c.jsx(B.StrictMode,{children:c.jsx(yy,{})}));

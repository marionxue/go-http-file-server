package asset

const mainJs = `
(function () {
var classNone = 'none';
var classHeader = 'header';
var leavingEvent = typeof window.onpagehide !== 'undefined' ? 'pagehide' : 'beforeunload';
var Enter = 'Enter';
var Escape = 'Escape';
var Esc = 'Esc';
var Space = ' ';
function enableFilter() {
if (!document.querySelector) {
var filter = document.getElementById && document.getElementById('panel-filter');
if (filter) {
filter.className += ' none';
}
return;
}
// pre check
var filter = document.body.querySelector('.filter');
if (!filter) {
return;
}
if (!filter.classList || !filter.addEventListener) {
filter.className += ' none';
return;
}
var input = filter.querySelector('input.filter-text');
if (!input) {
return;
}
var selectorNone = '.' + classNone;
var selectorNotNone = ':not(' + selectorNone + ')';
var selectorItem = '.item-list > li:not(.' + classHeader + '):not(.parent)';
var selectorItemNone = selectorItem + selectorNone;
var selectorItemNotNone = selectorItem + selectorNotNone;
// event handler
var timeoutId;
var lastFilterText = '';
var doFilter = function () {
var filterText = input.value.trim().toLowerCase();
if (filterText === lastFilterText) {
return;
}
var selector, items, i;
if (!filterText) {	// filter cleared, show all items
selector = selectorItemNone;
items = document.body.querySelectorAll(selector);
for (i = items.length - 1; i >= 0; i--) {
items[i].classList.remove(classNone);
}
} else {
if (filterText.indexOf(lastFilterText) >= 0) {	// increment search, find in visible items
selector = selectorItemNotNone;
} else if (lastFilterText.indexOf(filterText) >= 0) {	// decrement search, find in hidden items
selector = selectorItemNone;
} else {
selector = selectorItem;
}
items = document.body.querySelectorAll(selector);
for (i = items.length - 1; i >= 0; i--) {
var item = items[i];
var name = item.querySelector('.name');
if (name && name.textContent.toLowerCase().indexOf(filterText) < 0) {
item.classList.add(classNone);
} else {
item.classList.remove(classNone);
}
}
}
lastFilterText = filterText;
};
var onValueMayChange = function () {
clearTimeout(timeoutId);
timeoutId = setTimeout(doFilter, 350);
};
input.addEventListener('input', onValueMayChange, false);
input.addEventListener('change', onValueMayChange, false);
input.addEventListener('keydown', function (e) {
switch (e.key) {
case Enter:
clearTimeout(timeoutId);
input.blur();
doFilter();
e.preventDefault();
break;
case Escape:
case Esc:
clearTimeout(timeoutId);
input.value = '';
doFilter();
e.preventDefault();
break;
}
}, false);
// init
if (sessionStorage) {
var prevSessionFilter = sessionStorage.getItem(location.pathname);
sessionStorage.removeItem(location.pathname);
window.addEventListener(leavingEvent, function () {
if (input.value) {
sessionStorage.setItem(location.pathname, input.value);
}
}, false);
if (prevSessionFilter) {
input.value = prevSessionFilter;
}
}
if (input.value) {
doFilter();
}
}
function enableKeyboardNavigate() {
if (
!document.querySelector ||
!document.addEventListener ||
!document.body.classList ||
!document.body.parentElement
) {
return;
}
var pathList = document.body.querySelector('.path-list');
var itemList = document.body.querySelector('.item-list');
if (!pathList && !itemList) {
return;
}
function getFocusableSibling(container, isPrev, startA) {
if (!container) {
return
}
if (!startA) {
startA = container.querySelector(':focus');
}
var startLI = startA;
while (startLI && startLI.tagName !== 'LI') {
startLI = startLI.parentElement;
}
if (!startLI) {
if (isPrev) {
startLI = container.firstElementChild;
} else {
startLI = container.lastElementChild;
}
}
if (!startLI) {
return;
}
var siblingLI = startLI;
do {
if (isPrev) {
siblingLI = siblingLI.previousElementSibling;
if (!siblingLI) {
siblingLI = container.lastElementChild;
}
} else {
siblingLI = siblingLI.nextElementSibling;
if (!siblingLI) {
siblingLI = container.firstElementChild;
}
}
} while (siblingLI !== startLI && (
siblingLI.classList.contains(classNone) ||
siblingLI.classList.contains(classHeader)
));
if (siblingLI) {
var siblingA = siblingLI.querySelector('a');
return siblingA;
}
}
var selectorFirstAvailLi = 'li:not(.' + classNone + '):not(.' + classHeader + ')';
function getFirstFocusableSibling(container) {
var li = container.querySelector(selectorFirstAvailLi);
var a = li && li.querySelector('a');
return a;
}
function getLastFocusableSibling(container) {
var a = container.querySelector('li a');
a = getFocusableSibling(container, true, a);
return a;
}
function getMatchedFocusableSibling(container, isPrev, startA, buf, key) {
var skipRound = buf === key;
var matchKeyA;
var firstCheckA;
var secondCheckA;
var a = startA;
do {
if (skipRound) {
skipRound = false;
continue;
}
if (!a) {
continue;
}
// firstCheckA maybe a focused a that not belongs to the list
// secondCheckA must be in the list
if (!firstCheckA) {
firstCheckA = a;
} else if (firstCheckA === a) {
return;
} else if (!secondCheckA) {
secondCheckA = a;
} else if (secondCheckA === a) {
return;
}
var textContent = (a.querySelector('.name') || a).textContent.toLowerCase();
if (buf.length <= textContent.length && textContent.substring(0, buf.length) === buf) {
return a;
}
if (!matchKeyA && textContent[0] === key) {
matchKeyA = a;
}
} while (a = getFocusableSibling(container, isPrev, a));
return matchKeyA;
}
var UP = 'Up';
var DOWN = 'Down';
var LEFT = 'Left';
var RIGHT = 'Right';
var ARROW_UP = 'ArrowUp';
var ARROW_DOWN = 'ArrowDown';
var ARROW_LEFT = 'ArrowLeft';
var ARROW_RIGHT = 'ArrowRight';
var ARROW_UP_CODE = 38;
var ARROW_DOWN_CODE = 40;
var ARROW_LEFT_CODE = 37;
var ARROW_RIGHT_CODE = 39;
var SKIP_TAGS = ['INPUT', 'BUTTON', 'TEXTAREA'];
var PLATFORM = navigator.platform;
var IS_MAC_PLATFORM = PLATFORM.indexOf('Mac') >= 0 || PLATFORM.indexOf('iPhone') >= 0 || PLATFORM.indexOf('iPad') >= 0 || PLATFORM.indexOf('iPod') >= 0
var lookupKey = '';
var lookupBuffer = '';
var lookupStartA = null;
var lookupTimer;
function delayClearLookupContext() {
clearTimeout(lookupTimer);
lookupTimer = setTimeout(function () {
lookupBuffer = '';
lookupStartA = null;
}, 850);
}
function lookup(key) {
key = key.toLowerCase();
if (key === lookupKey && key === lookupBuffer) {
// same as last key, lookup next for the same key as prefix
lookupStartA = itemList.querySelector(':focus');
lookupBuffer = lookupKey;
} else {
if (!lookupStartA) {
lookupStartA = itemList.querySelector(':focus');
}
lookupKey = key;
lookupBuffer += key;
}
delayClearLookupContext();
return getMatchedFocusableSibling(itemList, false, lookupStartA, lookupBuffer, key);
}
function getFocusItemByKeyPress(e) {
if (SKIP_TAGS.indexOf(e.target.tagName) >= 0) {
return;
}
if (e.key) {
if (!e.altKey && !e.shiftKey) {
switch (e.key) {
case LEFT:
case ARROW_LEFT:
if (e.ctrlKey || e.metaKey) {
return getFirstFocusableSibling(pathList);
} else {
return getFocusableSibling(pathList, true);
}
case RIGHT:
case ARROW_RIGHT:
if (e.ctrlKey || e.metaKey) {
return getLastFocusableSibling(pathList);
} else {
return getFocusableSibling(pathList, false);
}
case UP:
case ARROW_UP:
if (e.ctrlKey || e.metaKey) {
return getFirstFocusableSibling(itemList);
} else {
return getFocusableSibling(itemList, true);
}
case DOWN:
case ARROW_DOWN:
if (e.ctrlKey || e.metaKey) {
return getLastFocusableSibling(itemList);
} else {
return getFocusableSibling(itemList, false);
}
}
}
if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.key.length === 1) {
return lookup(e.key);
}
} else if (e.keyCode) {
if (!e.altKey && !e.shiftKey) {
switch (e.keyCode) {
case ARROW_LEFT_CODE:
if (e.ctrlKey || e.metaKey) {
return getFirstFocusableSibling(pathList);
} else {
return getFocusableSibling(pathList, true);
}
case ARROW_RIGHT_CODE:
if (e.ctrlKey || e.metaKey) {
return getLastFocusableSibling(pathList);
} else {
return getFocusableSibling(pathList, false);
}
case ARROW_UP_CODE:
if (e.ctrlKey || e.metaKey) {
return getFirstFocusableSibling(itemList);
} else {
return getFocusableSibling(itemList, true);
}
case ARROW_DOWN_CODE:
if (e.ctrlKey || e.metaKey) {
return getLastFocusableSibling(itemList);
} else {
return getFocusableSibling(itemList, false);
}
}
}
if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.keyCode >= 32 && e.keyCode <= 126) {
return lookup(String.fromCharCode(e.keyCode));
}
}
}
document.addEventListener('keydown', function (e) {
var newFocusEl = getFocusItemByKeyPress(e);
if (newFocusEl) {
e.preventDefault();
newFocusEl.focus();
}
});
}
function enhanceUpload() {
if (!document.querySelector || !document.addEventListener || !document.body.classList) {
return;
}
var upload = document.body.querySelector('.upload');
if (!upload) {
return;
}
var form = upload.querySelector('form');
if (!form) {
return;
}
var fileInput = form.querySelector('.file');
if (!fileInput) {
return;
}
var uploadType = document.body.querySelector('.upload-type');
if (!uploadType) {
return;
}
var file = 'file';
var dirFile = 'dirfile';
var innerDirFile = 'innerdirfile';
var optFile = uploadType.querySelector('.' + file);
var optDirFile = uploadType.querySelector('.' + dirFile);
var optInnerDirFile = uploadType.querySelector('.' + innerDirFile);
var optActive = optFile;
function addClass(ele, className) {
ele && ele.classList.add(className);
}
function removeClass(ele, className) {
ele && ele.classList.remove(className);
}
function enableAddDir() {
var classHidden = 'hidden';
var classActive = 'active';
function onClickOpt(optTarget, clearInput) {
if (optTarget === optActive) {
return;
}
removeClass(optActive, classActive);
optActive = optTarget;
addClass(optActive, classActive);
if (clearInput) {
fileInput.value = '';
}
return true;
}
function onClickOptFile(e) {
if (onClickOpt(optFile, Boolean(e))) {
fileInput.name = file;
fileInput.webkitdirectory = false;
}
}
function onClickOptDirFile(e) {
if (onClickOpt(optDirFile, Boolean(e))) {
fileInput.name = dirFile;
fileInput.webkitdirectory = true;
}
}
function onClickOptInnerDirFile(e) {
if (onClickOpt(optInnerDirFile, Boolean(e))) {
fileInput.name = innerDirFile;
fileInput.webkitdirectory = true;
}
}
function onKeydownOpt(e) {
switch (e.key) {
case Enter:
case Space:
if (e.ctrlKey || e.altKey || e.metaKey || e.shiftKey) {
break;
}
e.preventDefault();
e.stopPropagation();
if (e.target === optActive) {
break;
}
e.target.click();
break;
}
}
if (typeof fileInput.webkitdirectory === 'undefined') {
addClass(uploadType, classNone);
return;
}
optDirFile && removeClass(optDirFile, classHidden);
optInnerDirFile && removeClass(optInnerDirFile, classHidden);
if (optFile) {
optFile.addEventListener('click', onClickOptFile);
optFile.addEventListener('keydown', onKeydownOpt);
}
if (optDirFile) {
optDirFile.addEventListener('click', onClickOptDirFile);
optDirFile.addEventListener('keydown', onKeydownOpt);
}
if (optInnerDirFile) {
optInnerDirFile.addEventListener('click', onClickOptInnerDirFile);
optInnerDirFile.addEventListener('keydown', onKeydownOpt);
}
if (sessionStorage) {
var uploadTypeField = 'upload-type';
var prevUploadType = sessionStorage.getItem(uploadTypeField);
sessionStorage.removeItem(uploadTypeField);
window.addEventListener(leavingEvent, function () {
var activeUploadType = fileInput.name;
if (activeUploadType !== file) {
sessionStorage.setItem(uploadTypeField, activeUploadType)
}
}, false);
if (prevUploadType === dirFile) {
optDirFile && optDirFile.click();
} else if (prevUploadType === innerDirFile) {
optInnerDirFile && optInnerDirFile.click();
}
}
optFile && fileInput.addEventListener('change', function (e) {
// workaround fix for mobile device, select dir not work but still act like select files
// switch back to file
if (optActive === optFile) {
return;
}
var files = e.target.files;
if (!files || !files.length) {
return;
}
var nodir = Array.prototype.slice.call(files).every(function (file) {
return !file.webkitRelativePath;
});
if (nodir) {
onClickOptFile();	// prevent clear input files
}
});
}
function enableAddDragDrop() {
function onDragEnterOver(e) {
e.stopPropagation();
e.preventDefault();
addClass(e.currentTarget, 'dragging');
}
function onDragLeave(e) {
if (e.target === e.currentTarget) {
removeClass(e.currentTarget, 'dragging');
}
}
function onDrop(e) {
e.stopPropagation();
e.preventDefault();
removeClass(e.currentTarget, 'dragging');
fileInput.value = '';
if (!e.dataTransfer || !e.dataTransfer.files || !e.dataTransfer.files.length) {
return;
}
var items = Array.prototype.slice.call(e.dataTransfer.items);
if (items && items.length && items[0].webkitGetAsEntry) {
for (var i = 0, len = items.length; i < len; i++) {
var entry = items[i].webkitGetAsEntry();
if (entry && entry.isDirectory) {
return;
}
}
}
if (optFile && optActive !== optFile) {
optFile.focus();
optFile.click();
}
fileInput.files = e.dataTransfer.files;
}
upload.addEventListener('dragenter', onDragEnterOver, false);
upload.addEventListener('dragover', onDragEnterOver, false);
upload.addEventListener('dragleave', onDragLeave, false);
upload.addEventListener('drop', onDrop, false);
}
function enableUploadProgress() {	// also fix Safari upload filename has no path info
if (!FormData) {
return;
}
var btnSubmit = form.querySelector('.submit');
if (!btnSubmit) {
return;
}
function onComplete() {
btnSubmit.disabled = false;
}
function onLoad() {
location.reload();
}
form.addEventListener('submit', function (e) {
e.stopPropagation();
e.preventDefault();
var files = Array.prototype.slice.call(fileInput.files);
if (!files.length) {
return;
}
var formName = fileInput.name;
var parts = new FormData();
files.forEach(function (file) {
if (file.webkitRelativePath) {
parts.append(formName, file, file.webkitRelativePath);
} else {
parts.append(formName, file);
}
});
var xhr = new XMLHttpRequest();
xhr.addEventListener('error', onComplete);
xhr.addEventListener('load', onComplete);
xhr.addEventListener('load', onLoad);
xhr.open(form.method, form.action);
xhr.send(parts);
btnSubmit.disabled = true;
});
}
enableAddDir();
enableAddDragDrop();
enableUploadProgress();
}
function enableNonRefreshDelete() {
if (!document.querySelector) {
return;
}
var itemList = document.body.querySelector('.item-list');
if (!itemList || !itemList.addEventListener) {
return;
}
if (itemList.classList) {
if (!itemList.classList.contains('has-deletable')) {
return;
}
} else if (itemList.className.indexOf('has-deletable') < 0) {
return;
}
itemList.addEventListener('click', function (e) {
if (e.defaultPrevented || !e.target || e.target.className.indexOf('delete') < 0) {
return;
}
var target = e.target;
var xhr = new XMLHttpRequest();
xhr.open('POST', target.href);
xhr.onload = function () {
var item = target;
var parentNode = item.parentNode;
while (item.nodeName !== 'LI') {
if (!parentNode) {
break;
}
item = parentNode;
parentNode = item.parentNode;
}
if (parentNode) {
parentNode.removeChild(item);
}
item = null;
parentNode = null;
target = null;
};
xhr.onerror = xhr.onabort = function () {
target = null;
};
xhr.send();
e.preventDefault();
return false;
}, false);
}
enableFilter();
enableKeyboardNavigate();
enhanceUpload();
enableNonRefreshDelete();
})();
`

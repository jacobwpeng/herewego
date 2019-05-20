var WALL_PROXY = '{{.WallProxy}}';

var NO_WALL_PROXY = '{{.NoWallProxy}}';

var DIRECT = '{{.Direct}}';

var NO_WALL_PROXY_SITES = [
{{ range $key, $value := .NoWallSites -}}
  "{{$key}}",
{{end}}
];

var DIRECT_SITES = [
{{ range $key, $value := .DirectSites -}}
  "{{$key}}",
{{end}}
];

function FindProxyForURL(url, host) {
    function check_private_ipaddr() {
      var re = /^(?:10|127|172\.(?:1[6-9]|2[0-9]|3[01])|192\.168)\..*/g;
      if (re.test(host)) {
        return true
      }
    }

    function check_ipv4() {
      var re_ipv4 = /^\d+\.\d+\.\d+\.\d+$/g;
      if (re_ipv4.test(host)) {
          return true;
      }
    }

    function isDomain(domain) {
        return ((domain[0] === '.') ? (host === domain.slice(1) || ((host_length = host.length) >= (domain_length = domain.length) && host.slice(host_length - domain_length) === domain)) : (host === domain));
    }

    function no_wall_proxy_site(callback) {
      for (var i = 0; i < NO_WALL_PROXY_SITES.length; i++) {
        if (callback(NO_WALL_PROXY_SITES[i]) === true) {
          return true;
        }
      }
      return false;
    }

    function direct_site(callback) {
      for (var i = 0; i < DIRECT_SITES.length; i++) {
        if (callback(DIRECT_SITES[i]) === true) {
          return true;
        }
      }
      return false;
    }

    if (direct_site(isDomain) || check_private_ipaddr()) {
      return DIRECT;
    } else if (isPlainHostName(host) === true 
            || check_ipv4() === true 
            || no_wall_proxy_site(isDomain) === true) {
      return NO_WALL_PROXY;
    } else {
      return WALL_PROXY;
    }

}
/*
    MIT License
    Copyright (C) 2012 n0gfwall0@gmail.com

    Permission is hereby granted, free of charge, to any person obtaining a 
    copy of this software and associated documentation files (the "Software"), 
    to deal in the Software without restriction, including without limitation 
    the rights to use, copy, modify, merge, publish, distribute, sublicense, 
    and/or sell copies of the Software, and to permit persons to whom the 
    Software is furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in 
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR 
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, 
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE 
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER 
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING 
    FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
    IN THE SOFTWARE.

                                                                              */

// SPDX-License-Identifier: MIT
pragma solidity >=0.5.0 <0.9.0;

address constant ADDRESS = {{ .Addr }};

library console {
    function _log(bytes memory data) internal view {
        assembly {
            pop(staticcall(gas(), ADDRESS, add(data, 0x20), mload(data), 0, 0))
        }
    }

    {{- range .Signatures }}

    function log(
        {{- range $index, $type := . -}}
            {{- if $index }}, {{ end -}}
            {{ $type | typeToArg }} arg{{ $index }}
        {{- end -}}
    ) internal view {
        _log(abi.encodeWithSignature("log(
            {{- range $index, $type := . -}}
                {{- if $index }},{{ end -}}
                {{- $type -}}
            {{- end -}}
        )", {{ range $index, $type := . -}}
                {{- if $index }}, {{ end -}}
                arg{{ $index -}}
            {{- end -}}
        ));
    }
    {{- end }}
}

library state {
    function dumpMemory() internal view {
        assembly {
            pop(staticcall(gas(), ADDRESS, 0, mload(0x40), 0, 0))
        }
    }

    function dumpMemory(uint start, uint length) internal view {
        assembly {
            pop(staticcall(gas(), ADDRESS, start, length, 0, 0))
        }
    }
}

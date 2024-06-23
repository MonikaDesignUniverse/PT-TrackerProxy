<?php

namespace App\Helpers;

use App\Models\User;
use Illuminate\Http\Request;

class PTTPHelper
{
    private static function ProcessPTTPListen(?string $pttpAddr, int|string|null $pttpPort): array {
        if ($pttpAddr == null || !filter_var($pttpAddr, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
            $pttpAddr = '127.0.0.1';
        }
        $pttpPort = null;
        if ($pttpPort == null || !is_numeric($pttpPort) || $pttpPort <= 1024 || $pttpPort >= 65535) {
            $pttpPort = 7887;
        }

        return array($pttpAddr, $pttpPort);
    }
    public static function GetPTTPInfo(Request $request): array {
        $pttpVersion = $request->headers->get('x-pttp-version');

        if ($pttpVersion != null) {
            /*
            $pttpPort = $request->headers->get('x-pttp-listenport');
            */
            list($pttpAddr, $pttpPort) = self::ProcessPTTPListen($request->headers->get('x-pttp-listenaddr'), null);
            return array('usePTTP' => true, 'pttpAddr' => $pttpAddr, 'pttpPort' => $pttpPort);
        }

        return array('usePTTP' => false);
    }
    public static function GetPTTPInfoFromUser(User $user): array {
        if (!empty($user->pttp_listen)) {
            $pttpListenSplit = explode(':', $user->pttp_listen, 3);
            if (count($pttpListenSplit) === 2) {
                list($pttpAddr, $pttpPort) = self::ProcessPTTPListen($pttpListenSplit[0], $pttpListenSplit[1]);
                return array('usePTTP' => true, 'pttpAddr' => $pttpAddr, 'pttpPort' => $pttpPort);
            }
        }

        return array('usePTTP' => false);
    }
    public static function GetPTTPURL(array $pttpInfo, string $url) {
        if ($pttpInfo['usePTTP']) {
            return str_replace(array('http://', 'https://'), "http://{$pttpInfo['pttpAddr']}:{$pttpInfo['pttpPort']}/", $url);
        }

        return $url;
    }
    public static function GetPTTPIPs(Request $request): array {
        $pttpIPs = [];

        $pttpIP4 = $request->headers->get('x-pttp-ip4');
        $pttpIP6 = $request->headers->get('x-pttp-ip6');
        if ($pttpIP4 != null && filter_var($pttpIP4, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
            array_push($pttpIPs, $pttpIP4);
        }
        if ($pttpIP6 != null && filter_var($pttpIP6, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6 | FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
            array_push($pttpIPs, $pttpIP6);
        }

        return $pttpIPs;
    }
}

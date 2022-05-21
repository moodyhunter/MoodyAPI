import json
import os
import utils.printer as P
from subprocess import check_output


class NetworkConfig:
    name: str
    ipv4_subnet: str
    ipv4_subnet_prefix: int

    def __init__(self, network_conf: dict) -> None:
        self.name = network_conf["name"]
        self.ipv4_subnet = network_conf["ipv4_subnet"]
        self.ipv4_subnet_prefix = network_conf["ipv4_subnet_prefix"]


class NetworkNode:
    name: str
    isServer: bool
    ipv4: str
    ipv6: str
    port: int
    public_key: str
    private_key: str
    extra_allowed_ips: list[str]
    public_endpoint: str

    def __init__(self, client_conf: dict) -> None:
        self.name = client_conf["name"]
        self.isServer = client_conf.get("isServer", False)
        self.ipv4 = client_conf.get("ipv4")
        self.ipv6 = client_conf.get("ipv6", "")
        self.port = client_conf.get("port", 51820)
        self.public_key = client_conf.get("public_key", "")
        self.private_key = client_conf.get("private_key", "")
        self.extra_allowed_ips = client_conf.get("extra_allowed_ips", [])
        self.public_endpoint = client_conf.get("public_endpoint", "")

    def generate_config(self) -> None:
        pass


class Network:
    name: str
    network_config: NetworkConfig
    nodes: dict[str, NetworkNode]

    def __init__(self, network_config: NetworkConfig, nodes: dict[str, NetworkNode]) -> None:
        self.name = network_config.name
        self.network_config = network_config
        self.nodes = nodes

        for _, node in self.nodes.items():
            node.public_key = check_output(["wg", 'pubkey'], input=node.private_key, text=True).strip()

    def generate_node_configuration(self, node_name: str) -> None:
        node = self.nodes[node_name]
        if node is None:
            return ""

        buffer = ""
        buffer += "[Interface]"

        buffer += "\n"

        if node.ipv4 != "":
            buffer += "Address = {}".format(node.ipv4)
            # TODO is this correct?
            if node.isServer:
                buffer += "/{}".format(self.network_config.ipv4_subnet_prefix)
            else:
                buffer += "/32"
            buffer += "\n"

        if node.ipv6 != "":
            buffer += "Address = {}\n".format(node.ipv6)

        buffer += "ListenPort = {}\n".format(node.port)
        buffer += "PrivateKey = {}\n".format(node.private_key)

        for peer_name, peer_node in self.nodes.items():
            # Don't connect to self
            if peer_name == node_name:
                continue

            # If the peer's endpoint is empty...
            if peer_node.public_endpoint == "":
                # and if it's a server, ...then we can't connect to it, so fail
                if peer_node.isServer:
                    P.error("Network node '{}' is a server but has no public endpoint.", [peer_name])

                # it's a client AND we are a client, then skip this client peer
                if not node.isServer:
                    P.sub_warning("Network node '{}' is a client with no public endpoint, since we are not a server, skipping.", [peer_name])
                    continue

            buffer += "\n"
            buffer += "# {} '{}'\n".format("Server" if peer_node.isServer else "Client", peer_name)
            buffer += "[Peer]\n"
            buffer += "PublicKey = {}\n".format(peer_node.public_key)

            allowed_ips = peer_node.extra_allowed_ips.copy()

            if peer_node.isServer:
                allowed_ips.append(self.network_config.ipv4_subnet + "/" + str(self.network_config.ipv4_subnet_prefix))

            if node.isServer:
                allowed_ips.append(peer_node.ipv4 + "/32")
                if peer_node.ipv6 != "":
                    allowed_ips.append(peer_node.ipv6 + "/128")

            if len(allowed_ips) > 0:
                buffer += "AllowedIPs = {}\n".format(", ".join(set(allowed_ips)))

            if peer_node.public_endpoint != "":
                buffer += "Endpoint = {}\n".format(peer_node.public_endpoint)

            buffer += "PersistentKeepalive = 10\n"

        print(buffer)


def load_network(network_name: str) -> Network:
    P.progress("Reading configurations from directory '{}':", ["networks/" + network_name])

    files = os.listdir('networks/' + network_name)
    files.sort()

    nodes: dict[str, NetworkNode] = {}
    network_conf: NetworkConfig = None

    for entry in files:
        full_path = os.path.join('networks/' + network_name, entry)

        try:
            if full_path.endswith("/_network_.json"):
                with open(full_path, 'r') as f:
                    data = json.load(f)
                    network_conf = NetworkConfig(data)
                    P.sub_progress("Network Name: {}", [network_conf.name])
            else:
                with open(full_path) as f:
                    data = json.load(f)
                    n = NetworkNode(data)
                    nodes[n.name] = n
                    P.sub_progress("Server: {}" if n.isServer else "Client: {}", [n.name])

        except Exception as e:
            P.error("Failed to read '{}': {}", [entry, e])

    if len(nodes) == 0:
        P.error("No nodes available for network'{}'", [network_name])

    if network_conf is None:
        P.error("No network configuration found for network '{}'", [network_name])

    P.progress("Network configuration loaded:")
    P.sub_progress("Network Name: {}", [network_conf.name])
    P.sub_progress("IPv4 Subnet: {}", [network_conf.ipv4_subnet + "/" + str(network_conf.ipv4_subnet_prefix)])
    P.sub_progress("Found {} Nodes", [len(nodes)])

    return Network(network_conf, nodes)

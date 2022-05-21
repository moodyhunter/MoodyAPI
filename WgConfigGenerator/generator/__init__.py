import json
import os
import utils.printer as P


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
        self.port = client_conf.get("port", 0)
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

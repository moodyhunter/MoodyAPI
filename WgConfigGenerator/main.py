from genericpath import isdir
import os
import enquiries
import generator
import utils.printer as P


def main():
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    network_list = []

    for entry in os.listdir('networks'):
        if isdir(os.path.join('networks', entry)):
            network_list.append(entry)

    if len(network_list) == 0:
        P.error("No networks found in '{}'", "networks/")

    P.progress("Found {} network(s): {}", [len(network_list), network_list])

    if len(network_list) > 1 and "demo_network" in network_list:
        network_list.remove("demo_network")
        P.sub_progress("Removing the demo network '{}' from list", ["demo_network"])

    if len(network_list) == 1:
        chosen_networks = network_list
        P.sub_progress("Only one network available, using it directly")
    else:
        chosen_networks = enquiries.choose('Select the network(s) you want to configure:', network_list, multi=True)

    if len(chosen_networks) == 0:
        P.error("No networks selected")

    for network_name in chosen_networks:
        network = generator.load_network(network_name)

        nodes = network.nodes.keys()
        chosen_nodes = enquiries.choose('Select the node(s) you want to generate:', nodes, multi=True)

        P.progress("Generating configuration for network '{}'", [network_name])
        P.sub_progress("Generating configuration for {} node(s)", [len(chosen_nodes)])

        for node_name in chosen_nodes:
            node = network.nodes[node_name]
            node.generate_config()
            P.sub_progress("Generated configuration for node '{}'", [node_name])

    pass


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nExiting...")

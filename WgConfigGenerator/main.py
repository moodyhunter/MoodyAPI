#!/usr/bin/env python
from genericpath import isdir
import os
import shutil
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

    P.progress("Found {} network(s): {}", [len(network_list), ", ".join(network_list)])

    if len(network_list) > 1 and "demo_network" in network_list:
        network_list.remove("demo_network")
        P.sub_progress("Removing the demo network '{}' from list", ["demo_network"])

    if len(network_list) == 1:
        chosen_network_dir = network_list[0]
        P.sub_progress("Only one network available, using it directly")
    else:
        P.sub_progress("Select the network you want to configure:")
        chosen_network_dir = enquiries.choose('', network_list)[0]

    network = generator.load_network(chosen_network_dir)
    all_nodes = network.nodes.keys()
    P.progress("Select the node(s) you want to generate (empty selection means select all):")
    chosen_nodes = enquiries.choose('', all_nodes, multi=True)

    if len(chosen_nodes) == 0:
        P.sub_warning("No nodes selected, using all nodes")
        chosen_nodes = all_nodes

    P.progress("Checking for existing configurations...", [])

    generated_dir = os.path.join('generated', chosen_network_dir)

    os.makedirs(generated_dir, exist_ok=True)
    existing_configurations = os.listdir(generated_dir)

    if len(existing_configurations) > 0:
        P.sub_progress("Found {} existing configuration(s): {}", [len(existing_configurations), ", ".join(existing_configurations)])
        should_remove = enquiries.confirm("Do you want to remove ALL of them?")
        if should_remove:
            shutil.rmtree(generated_dir)
            P.sub_progress("Removed all existing configurations for network '{}'", [network.name])
            os.makedirs(generated_dir, exist_ok=True)
        else:
            P.sub_warning("Not removing existing configurations, some configurations may be outdated")

    P.progress("Generating configurations...", [])

    for node_name in chosen_nodes:
        P.sub_progress("Generated configuration for node '{}'", [node_name])
        text = network.generate_node_configuration(node_name)
        print(text)

        import pyqrcode
        print(pyqrcode.create(text).terminal())

    pass


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nCtrl+C pressed, exiting...")

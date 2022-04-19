/*
Copyright © 2022 Nethermind hello.nethermind.io

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/NethermindEth/1click/configs"
	"github.com/NethermindEth/1click/internal/pkg/clients"
	"github.com/NethermindEth/1click/internal/pkg/generate"
	"github.com/NethermindEth/1click/internal/ui"
	"github.com/NethermindEth/1click/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	executionName  string
	consensusName  string
	validatorName  string
	generationPath string
	randomize      bool
	install        bool
	run            bool
	y              bool
	services       *[]string
)

const (
	execution, consensus, validator = "execution", "consensus", "validator"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli [flags]",
	Short: "Quick start 1click",
	Long: `Run the setup tool on-premise in a quick way. Provide only the command line
options and the tool will do all the work.

First it will check if dependencies like docker and docker-compose are installed on your machine
and provide instructions for installing them if they are not installed.

Second, it will generate docker-compose scripts to run the full setup according to your selection.

Finally, it will run the generated docker-compose script. Only execution and consensus clients will be executed by default.

Running the command without flags (except global flag'--config') is equivalent to '1click cli -r' `,
	Args: cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Count flags being set
		count := 0
		// HACKME: LocalFlags() doesn't work, so we count manually and check for parent flag config
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if f.Name != "config" {
				count++
			}
		})

		if count == 0 {
			// No flag behavior
			randomize = true
		}

		// Quick run
		if y {
			randomize, install, run = true, true, true
		}

		// Validate run-clients flag
		if utils.Contains(*services, "all") {
			if len(*services) == 1 {
				// all used correctly
				services = &[]string{execution, consensus, validator}
			} else {
				// Ambiguous value
				log.Fatalf(configs.RunClientsFlagAmbiguousError, *services)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get all clients: supported + configured
		clientsMap, errors := clients.GetClients([]string{execution, consensus, validator})
		if len(errors) > 0 {
			for _, err := range errors {
				log.Error(err)
			}
			os.Exit(1)
		}

		// Handle selection and validation of clients
		combinedClients, err := validateClients(clientsMap)
		if err != nil {
			log.Fatal(err)
		}

		dependencies := configs.GetDependencies()
		log.Infof(configs.CheckingDependencies, strings.Join(dependencies, ", "))

		// Check if dependencies are installed. Keep checking dependencies until they are all installed
		for pending := utils.CheckDependencies(dependencies); len(pending) > 0; pending = utils.CheckDependencies(dependencies) {
			log.Infof(configs.DependenciesPending, strings.Join(pending, ", "))
			if install {
				// Install dependencies directly
				if err := installDependencies(pending); err != nil {
					log.Fatal(err)
				}
			} else {
				// Let the user decide to see the instructions for installing dependencies and exit or let the tool install them and continue
				if err := installOrShowInstructions(pending); err != nil {
					log.Fatal(err)
				}
			}
		}
		log.Info(configs.DependenciesOK)

		// Generate docker-compose scripts
		if err = generate.GenerateScripts(combinedClients.Execution.Name, combinedClients.Consensus.Name, combinedClients.Validator.Name, generationPath); err != nil {
			log.Fatal(err)
		}

		if run {
			if err = runAndShowContainers(); err != nil {
				log.Fatal(err)
			}
		} else {
			// Let the user decide to see the instructions for executing the scripts and exit or let the tool execute them
			if err = runScriptOrExit(); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Local flags
	cliCmd.Flags().StringVarP(&executionName, "execution", "e", "", "Execution engine client, e.g. Geth, Nethermind, Besu, Erigon")

	cliCmd.Flags().StringVarP(&consensusName, "consensus", "c", "", "Consensus engine client, e.g. Teku, Lodestar, Prysm, Lighthouse, Nimbus")

	cliCmd.Flags().StringVarP(&validatorName, "validator", "v", "", "Validator engine client, e.g. Teku, Lodestar, Prysm, Lighthouse, Nimbus")

	cliCmd.Flags().StringVarP(&generationPath, "path", "p", configs.DefaultDockerComposeScriptsPath, "docker-compose scripts generation path")

	cliCmd.Flags().BoolVarP(&randomize, "randomize", "r", false, "Randomize combination of clients")

	cliCmd.Flags().BoolVarP(&install, "install", "i", false, "Install dependencies if not installed without asking")

	cliCmd.Flags().BoolVar(&run, "run", false, "Run the generated docker-compose scripts without asking")

	cliCmd.Flags().BoolVarP(&y, "yes", "y", false, "Shortcut for '1click cli -r -i --run'. Run without prompts")

	services = cliCmd.Flags().StringSlice("run-clients", []string{execution, consensus}, "Run only the specified clients. Possible values: execution, consensus, validator, all. The 'all' option must be used alone. Example: '1click cli -r --run-clients=consensus,validator'")
}

func installOrShowInstructions(pending []string) (err error) {
	optInstall, optExit := "Install dependencies", "Exit. You will manage this dependencies on your own"
	prompt := promptui.Select{
		Label: "Select how to proceed with the pending dependencies",
		Items: []string{optInstall, optExit},
	}

	if err = utils.HandleInstructions(pending, utils.ShowInstructions); err != nil {
		return fmt.Errorf(configs.ShowingInstructionsError, err)
	}
	_, result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf(configs.PromptFailedError, err)
	}

	switch result {
	case optInstall:
		return installDependencies(pending)
	default:
		log.Info(configs.Exiting)
		os.Exit(0)
	}

	return nil
}

func installDependencies(pending []string) error {
	if err := utils.HandleInstructions(pending, utils.InstallDependency); err != nil {
		return fmt.Errorf(configs.InstallingDependenciesError, err)
	}
	return nil
}

func randomizeClients(allClients clients.OrderedClients) (clients.Clients, error) {
	var executionClient, consensusClient, validatorClient clients.Client
	var combinedClients clients.Clients

	executionClient, err := clients.RandomChoice(allClients[execution])
	if err != nil {
		return combinedClients, err
	}
	consensusClient, err = clients.RandomChoice(allClients[consensus])
	if err != nil {
		return combinedClients, err
	}
	validatorClient, err = clients.RandomChoice(allClients[validator])
	if err != nil {
		return combinedClients, err
	}

	combinedClients = clients.Clients{
		Execution: executionClient,
		Consensus: consensusClient,
		Validator: validatorClient}
	return combinedClients, nil
}

func validateClients(allClients clients.OrderedClients) (clients.Clients, error) {
	var combinedClients clients.Clients
	var err error

	if randomize {
		// Select a random execution client, consensus client and validator client
		combinedClients, err = randomizeClients(allClients)
		if err != nil {
			return combinedClients, err
		}

		log.Infof("Listing randomized clients\n\n")
		ui.WriteRandomizedClientsTable([][]string{
			{"Execution", combinedClients.Execution.Name},
			{"Consensus", combinedClients.Consensus.Name},
			{"Validator", combinedClients.Validator.Name},
		})
	} else {
		notProvidedClients := make([]string, 0)
		if executionName == "" {
			notProvidedClients = append(notProvidedClients, execution+" client")
		}
		if consensusName == "" {
			notProvidedClients = append(notProvidedClients, consensus+" client")
		}
		if validatorName == "" {
			notProvidedClients = append(notProvidedClients, validator+" client")
		}

		if len(notProvidedClients) > 0 {
			var msg string

			if len(notProvidedClients) == 1 {
				msg = notProvidedClients[0]
			} else {
				msg = strings.Join(notProvidedClients[:len(notProvidedClients)-1], ", ")
				msg = msg + " and " + notProvidedClients[len(notProvidedClients)-1]
			}

			return combinedClients, fmt.Errorf(configs.ClientNotSpecifiedError, msg)
		}

		exec, ok := allClients[execution][executionName]
		if !ok {
			exec.Name = executionName
		}
		cons, ok := allClients[consensus][consensusName]
		if !ok {
			cons.Name = consensusName
		}
		val, ok := allClients[validator][validatorName]
		if !ok {
			val.Name = validatorName
		}

		combinedClients = clients.Clients{
			Execution: exec,
			Consensus: cons,
			Validator: val,
		}
	}

	if err = clients.ValidateClient(combinedClients.Execution, execution); err != nil {
		return combinedClients, err
	}
	if err = clients.ValidateClient(combinedClients.Consensus, consensus); err != nil {
		return combinedClients, err
	}
	if err = clients.ValidateClient(combinedClients.Validator, validator); err != nil {
		return combinedClients, err
	}

	return combinedClients, nil
}

func runScriptOrExit() (err error) {
	optRun, optExit := "Run the script", "Exit"
	prompt := promptui.Select{
		Label: "Select how to proceed with the generated docker-compose script",
		Items: []string{optRun, optExit},
	}

	log.Infof(configs.InstructionsFor, "running docker-compose script")
	servs := strings.Join(*services, " ")
	fmt.Printf("\n%s\n\n", fmt.Sprintf(configs.DockerComposeUpCMD, generationPath, servs))

	_, result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("prompt failed %s", err)
	}

	switch result {
	case optRun:
		if err = runAndShowContainers(); err != nil {
			return err
		}
	default:
		log.Info(configs.Exiting)
		os.Exit(0)
	}

	return nil
}

func runAndShowContainers() error {
	// TODO: (refac) Put this check to checks.go and call it from there
	// Check if docker engine is on
	log.Info(configs.CheckingDockerEngine)
	log.Infof(configs.RunningCommand, configs.DockerPsCMD)
	if _, err := utils.RunCmd(configs.DockerPsCMD, true, false); err != nil {
		return fmt.Errorf(configs.DockerEngineOffError, err)
	}

	// Run docker-compose script
	servs := strings.Join(*services, " ")
	upCMD := fmt.Sprintf(configs.DockerComposeUpCMD, generationPath+"/docker-compose.yml", servs)
	log.Infof(configs.RunningCommand, upCMD)
	if _, err := utils.RunCmd(upCMD, false, false); err != nil {
		return fmt.Errorf(configs.CommandError, upCMD, err)
	}

	// Run docker-compose ps --filter status=running to show script running containers
	psCMD := fmt.Sprintf(configs.DockerComposePsFilterCMD, generationPath+"/"+configs.DefaultDockerComposeScriptName)
	log.Infof(configs.RunningCommand, psCMD)
	if _, err := utils.RunCmd(psCMD, false, false); err != nil {
		return fmt.Errorf(configs.CommandError, psCMD, err)
	}

	return nil
}
